package ports

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alukart32/effective-mobile-test-task/internal/person/model"
	"github.com/alukart32/effective-mobile-test-task/internal/pkg/zerologx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func HttpRoutes(router *gin.Engine, manager personManager) error {
	if manager == nil {
		return fmt.Errorf("init HTTP routes: personManager is nil")
	}

	g := router.Group("/persons")
	{
		g.GET("/", collectPersons(manager))
		g.POST("/", createPerson(manager))
		g.GET("/:id", getPerson(manager))
		g.DELETE("/:id", deletePerson(manager))
		g.PATCH("/:id", updatePerson(manager))
	}

	return nil
}

type createPersonRequest struct {
	Name       string
	Surname    string
	Patronymic string
}

type createPersonResponse struct {
	Id string
}

func createPerson(creator personCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := zerologx.Get().With().Ctx(c.Request.Context()).Logger()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("port", "http").Str("op", "create person")
		})

		var reqData createPersonRequest
		err := c.ShouldBindJSON(&reqData)
		if err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusBadRequest,
				gin.H{"err": fmt.Errorf("create person: %w", err).Error()})
			return
		}

		fio, err := model.NewFIO(reqData.Name, reqData.Surname, reqData.Patronymic)
		if err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusBadRequest,
				gin.H{"err": fmt.Errorf("create person: %w", err).Error()})
			return
		}
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Object("params", fio)
		})
		logger.Info().Msg(">> create person")

		var personId string
		if personId, err = creator.CreateFrom(c.Request.Context(), fio); err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusInternalServerError,
				gin.H{"err": fmt.Errorf("create person: %w", err).Error()})
			return
		}
		logger.Info().Str("person_id", personId).Msg("<< create person")
		c.JSON(http.StatusOK, createPersonResponse{Id: personId})
	}
}

func getPerson(finder personFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		logger := zerologx.Get().With().Ctx(c.Request.Context()).Logger()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("port", "http").
				Str("op", "find person by id").
				Str("param", id)
		})

		if len(id) == 0 {
			msg := "invalid value for id: empty"
			logger.Error().Msg(msg)
			c.JSON(http.StatusBadRequest,
				gin.H{"err": msg})
			return
		}
		logger.Info().Msg(">> find person")

		person, err := finder.FindById(c.Request.Context(), id)
		if err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusInternalServerError,
				gin.H{"err": fmt.Errorf("get person: %w", err).Error()})
			return
		}
		logger.Info().Object("person", person).Msg("<< find person")

		respBody, err := json.Marshal(person)
		if err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusInternalServerError,
				gin.H{"err": fmt.Errorf("get person: %w", err).Error()})
			return
		}
		c.Data(http.StatusOK, "application/json; charset=utf-8", respBody)
	}
}

func collectPersons(collector personCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			limit  int
			offset int
			filter model.PersonFilter
			err    error
		)
		logger := zerologx.Get().With().Ctx(c.Request.Context()).Logger()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("port", "http").Str("op", "collect persons")
		})
		limitFromQuery := c.Query("limit")
		if limitFromQuery != "" {
			limit, err = strconv.Atoi(limitFromQuery)
			if err != nil || limit < -1 {
				msg := "invalid value for limit: " + limitFromQuery
				logger.Error().Msg(msg)
				c.JSON(http.StatusBadRequest,
					gin.H{"err": msg})
				return
			}
		}
		offsetFromQuery := c.Query("offset")
		if offsetFromQuery != "" {
			offset, err = strconv.Atoi(offsetFromQuery)
			if err != nil || offset < -1 {
				msg := "invalid value for offset: " + offsetFromQuery
				logger.Error().Msg(msg)
				c.JSON(http.StatusBadRequest,
					gin.H{"err": msg})
				return
			}
		}
		filters := c.QueryArray("filter")
		if len(filters) > 0 {
			filter, err = model.NewPersonFilter(filters)
			if err != nil {
				msg := "invalid value for filter"
				logger.Err(err).Msg(msg)
				c.JSON(http.StatusBadRequest,
					gin.H{"err": fmt.Errorf("%s: %w", msg, err).Error()})
				return
			}
		}

		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Dict("params", zerolog.Dict().
				Int("limit", limit).
				Int("offset", offset).
				Object("filters", filter),
			)
		})
		logger.Info().Msg(">> collect persons")

		persons, err := collector.Collect(c.Request.Context(), filter, limit, offset)
		if err != nil {
			err = fmt.Errorf("collect persons: %w", err)
			logger.Err(err).Send()
			c.JSON(http.StatusInternalServerError,
				gin.H{"err": err.Error()})
			return
		}
		logger.Info().Str("status", "ok").Msg("<< collect persons")

		if len(persons) == 0 {
			logger.Info().Str("res", "no content")
			c.Status(http.StatusNoContent)
			return
		}

		respBody, err := json.Marshal(persons)
		if err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusInternalServerError,
				gin.H{"err": fmt.Errorf("collect persons: %w", err).Error()})
			return
		}
		c.Data(http.StatusOK, "application/json; charset=utf-8", respBody)
	}
}

type updatePersonRequest struct {
	Nation string
	Gender string
	Age    int
}

func updatePerson(updater personUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		logger := zerologx.Get().With().Ctx(c.Request.Context()).Logger()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("port", "http").
				Str("op", "update person").
				Str("patam", id)
		})

		if len(id) == 0 {
			msg := "invalid value for id: empty"
			logger.Error().Msg(msg)
			c.JSON(http.StatusBadRequest,
				gin.H{"err": msg})
			return
		}

		var reqData updatePersonRequest
		err := c.ShouldBindJSON(&reqData)
		if err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusBadRequest,
				gin.H{"err": fmt.Errorf("update person: %w", err).Error()})
			return
		}
		metaData := model.PersonalMetaData(reqData)
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Object("param", metaData)
		})
		logger.Info().Msg(">> update person")

		err = updater.Update(c.Request.Context(), id, metaData)
		if err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusInternalServerError,
				gin.H{"err": fmt.Errorf("update person: %w", err).Error()})
			return
		}
		logger.Info().Str("status", "ok").Msg("<< update person")
		c.Status(http.StatusOK)
	}
}

func deletePerson(deleter personDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		logger := zerologx.Get().With().Ctx(c.Request.Context()).Logger()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("port", "http").
				Str("op", "delete person").
				Str("param", id)
		})
		logger.Info().Msg(">> delete person")

		if len(id) == 0 {
			logger.Error().Msg("invalid value for id: empty")
			c.JSON(http.StatusBadRequest,
				gin.H{"err": "invalid value for id: empty"})
			return
		}

		err := deleter.Delete(c.Request.Context(), id)
		if err != nil {
			logger.Err(err).Send()
			c.JSON(http.StatusInternalServerError,
				gin.H{"err": fmt.Errorf("delete person: %w", err).Error()})
			return
		}
		logger.Info().Str("status", "ok").Msg("<< delete person")
		c.Status(http.StatusOK)
	}
}
