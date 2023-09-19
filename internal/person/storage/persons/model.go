package persons

import (
	"encoding/json"

	"github.com/alukart32/effective-mobile-test-task/internal/person/model"
)

type record struct {
	Id         string `redis:"id" json:"id"`
	Name       string `redis:"name" json:"name"`
	Surname    string `redis:"surname" json:"surname"`
	Patronymic string `redis:"patronymic" json:"patronymic"`
	Nation     string `redis:"nation" json:"nation"`
	Gender     string `redis:"gender" json:"gender"`
	Age        int    `redis:"age" json:"age"`
}

func (r record) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func toRecord(p model.Person) record {
	return record{
		Id:         p.Id,
		Name:       p.Name,
		Surname:    p.Surname,
		Patronymic: p.Patronymic,
		Nation:     p.Nation,
		Gender:     p.Gender,
		Age:        p.Age,
	}
}

func (r record) ToModel() model.Person {
	return model.Person{
		Id: r.Id,
		FIO: model.FIO{
			Name:       r.Name,
			Surname:    r.Surname,
			Patronymic: r.Patronymic,
		},
		PersonalMetaData: model.PersonalMetaData{
			Nation: r.Nation,
			Gender: r.Gender,
			Age:    r.Age,
		},
	}
}
