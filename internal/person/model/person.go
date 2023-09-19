package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Person2 struct {
	ID string
	FIO
	PersonalMetaData
}

type Person struct {
	Id string
	FIO
	PersonalMetaData
}

func NewPerson(
	fio FIO,
	meta PersonalMetaData,
) Person {
	return Person{
		Id:               uuid.New().String(),
		FIO:              fio,
		PersonalMetaData: meta,
	}
}

func (p Person) IsEmpty() bool {
	return len(p.Name) == 0 && len(p.Surname) == 0 &&
		len(p.Nation) == 0 && len(p.Gender) == 0 &&
		p.Age == 0
}

func (p Person) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", p.Id).
		Object("fio", p.FIO).
		Object("meta", p.PersonalMetaData)
}

type PersonFilter struct {
	OlderThan, YoungerThan int
	Gender                 string
	Nations                []string
}

func NewPersonFilter(filters []string) (PersonFilter, error) {
	if len(filters) == 0 {
		return PersonFilter{}, nil
	}

	var filter PersonFilter
	for _, f := range filters {
		vals := strings.Split(f, ".")
		switch vals[0] {
		case "older-than":
			if len(vals[1]) != 0 {
				i, err := strconv.Atoi(vals[1])
				if err != nil {
					return PersonFilter{}, fmt.Errorf("filter parsing error: %s", f)
				}
				filter.OlderThan = i
			}
		case "younger-than":
			if len(vals[1]) != 0 {
				i, err := strconv.Atoi(vals[1])
				if err != nil {
					return PersonFilter{}, fmt.Errorf("filter parsing error: %s", f)
				}
				filter.YoungerThan = i
			}
		case "gender":
			if len(vals[1]) != 0 {
				filter.Gender = vals[1]
			}
		case "nation":
			if len(vals[1]) != 0 {
				filter.Nations = append(filter.Nations, vals[1])
			}
		default:
			return PersonFilter{}, fmt.Errorf("unsupported %s filter", f)
		}
	}
	return filter, nil
}

func (f PersonFilter) IsEmpty() bool {
	return f.YoungerThan == 0 && f.OlderThan == 0 &&
		len(f.Gender) == 0 && len(f.Nations) == 0
}

func (f PersonFilter) MarshalZerologObject(e *zerolog.Event) {
	e.
		Int("older-than", f.OlderThan).
		Int("younger-than", f.YoungerThan).
		Str("gender", f.Gender)

	nations := zerolog.Arr()
	for _, n := range f.Nations {
		nations.Str(n)
	}

	e.Array("nations", nations)
}

type PersonalMetaData struct {
	Nation string
	Gender string
	Age    int
}

func (meta PersonalMetaData) IsEmpty() bool {
	return len(meta.Nation) == 0 && len(meta.Gender) == 0 &&
		meta.Age == 0
}

func (meta PersonalMetaData) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("nation", meta.Nation).
		Str("gender", meta.Gender).
		Str("nation", meta.Nation)
}
