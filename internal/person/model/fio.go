package model

import (
	"fmt"
	"regexp"

	"github.com/rs/zerolog"
)

type FIO struct {
	Name       string
	Surname    string
	Patronymic string
}

var isAlpha = regexp.MustCompile(`^[A-Za-z]+$`).MatchString

func NewFIO(name string, surname string, patronymic string) (FIO, error) {
	var fio FIO

	if len(name) == 0 {
		return FIO{}, fmt.Errorf("empty required name")
	}
	if !isAlpha(name) {
		return FIO{}, fmt.Errorf("name contains invalid characters: %s", name)
	}
	fio.Name = name

	if len(surname) == 0 {
		return fio, fmt.Errorf("empty required surname")
	}
	if !isAlpha(surname) {
		return FIO{}, fmt.Errorf("surname contains invalid characters: %s", surname)
	}
	fio.Surname = surname

	if len(patronymic) != 0 && !isAlpha(patronymic) {
		return FIO{}, fmt.Errorf("patronymic contains invalid characters: %s", patronymic)
	}
	fio.Patronymic = patronymic
	return fio, nil
}

func (f FIO) String() string {
	return fmt.Sprintf("[name: %s, surname: %s, patronymic: %s]",
		f.Name, f.Surname, f.Patronymic)
}

func (f FIO) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("name", f.Name).
		Str("surname", f.Surname).
		Str("patronymic", f.Patronymic)
}
