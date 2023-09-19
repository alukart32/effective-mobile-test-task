package ports

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alukart32/effective-mobile-test-task/internal/person/model"
	"github.com/alukart32/effective-mobile-test-task/internal/pkg/zerologx"
	"github.com/rs/zerolog"
	kafka "github.com/segmentio/kafka-go"
)

type fioMsg struct {
	Name       string
	Surname    string
	Patronymic string
}

func (f fioMsg) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("name", f.Name).
		Str("surname", f.Surname).
		Str("patronymic", f.Patronymic)
}

func (m fioMsg) String() string {
	return fmt.Sprintf("[name: %s, surname: %s, patronymic: %s]",
		m.Name, m.Surname, m.Patronymic)
}

type fioErrorMsg struct {
	fioMsg
	err error
}

type kafkaFIO struct {
	reader      *kafka.Reader
	errorWriter *kafka.Writer

	msgs   chan fioMsg
	errors chan fioErrorMsg
	done   chan struct{}

	personCreator personCreator
}

func KafkaFIO(
	ctx context.Context,
	readTopic string,
	errorTopic string,
	brokers []string,
	bufferSize int,
	creator personCreator) (*kafkaFIO, error) {
	if len(readTopic) == 0 {
		return nil, fmt.Errorf("empty read topic")
	}
	if len(errorTopic) == 0 {
		return nil, fmt.Errorf("empty error topic")
	}
	if len(brokers) == 0 {
		return nil, fmt.Errorf("empty brokers list")
	}
	if bufferSize <= 0 {
		bufferSize = 1
	}
	if creator == nil {
		return nil, fmt.Errorf("person creator is nil")
	}

	logger := zerologx.Get().With().Logger()
	handler := kafkaFIO{
		personCreator: creator,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       readTopic,
			MaxBytes:    1e6,
			ErrorLogger: &logger,
		}),
		errorWriter: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  errorTopic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: false,
		},
		msgs:   make(chan fioMsg, bufferSize),
		errors: make(chan fioErrorMsg, 1),
		done:   make(chan struct{}, 1),
	}
	go handler.Handle(ctx)
	go handler.RespondError(ctx)
	go handler.Fetch(ctx)

	return &handler, nil
}

func (h *kafkaFIO) Handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-h.msgs:
			if !ok {
				return
			}

			fio, err := model.NewFIO(msg.Name, msg.Surname, msg.Patronymic)
			if err != nil {
				h.errors <- fioErrorMsg{fioMsg: msg, err: err}
				continue
			}

			_, err = h.personCreator.CreateFrom(ctx, fio)
			if err != nil {
				h.errors <- fioErrorMsg{fioMsg: msg, err: err}
			}
		}
	}
}

// Fetch reads msg FIO from kafka topic.
func (h *kafkaFIO) Fetch(ctx context.Context) {
	defer func() {
		close(h.msgs)
		if err := h.reader.Close(); err != nil {
			log.Fatal("kafka Fetch:", err)
		}
	}()
	logger := zerologx.Get().
		With().
		Str("port", "kafka").
		Logger()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		m, err := h.reader.ReadMessage(ctx)
		if err != nil {
			logger.Err(err).Send()
			break
		}
		logger.Info().Str("op", "fetch msg").Dict("msg",
			zerolog.Dict().
				Str("topic", m.Topic).
				Time("time", m.Time).
				Int("partion", m.Partition).
				Int64("offset", m.Offset).
				Str("key", string(m.Key)).
				Str("value", string(m.Value)),
		).Send()

		var msg fioMsg
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			logger.Err(err).Send()
			h.errors <- fioErrorMsg{err: err}
		} else {
			h.msgs <- msg
		}
	}
}

func (h *kafkaFIO) RespondError(ctx context.Context) {
	const retries = 3

	defer func() {
		if err := h.errorWriter.Close(); err != nil {
			log.Fatal("kafka RespondError:", err)
		}
	}()
	logger := zerologx.Get().
		With().
		Str("port", "kafka").
		Logger()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-h.errors:
			if !ok {
				return
			}
			logger.Info().
				Object("msg", msg).
				Msg("send msg")

			var (
				b   []byte
				err error
			)
			if b, err = json.Marshal(&msg); err != nil {
				logger.Err(err).Send()
				continue
			}

			messages := []kafka.Message{
				{
					Value: b,
				},
			}
			for i := 0; i < retries; i++ {
				select {
				case <-h.done:
					return
				default:
				}

				writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				err = h.errorWriter.WriteMessages(writeCtx, messages...)
				if errors.Is(err, kafka.LeaderNotAvailable) ||
					errors.Is(err, context.DeadlineExceeded) {
					logger.Info().Err(err).Send()
					<-time.After(time.Millisecond * 250)
					continue
				}
				if err != nil {
					logger.Err(err).Send()
				}
				break
			}
		}
	}
}
