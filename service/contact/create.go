package contact

import (
	contactRepositoryDto "github.com/mohamadrezamomeni/graph/dto/repository/contact"
	contactServiceDto "github.com/mohamadrezamomeni/graph/dto/service/contact"
)

type Contact struct {
	contactRepo ContactRepo
}

type ContactRepo interface {
	Create(*contactRepositoryDto.Create) error
}

func New(contactRepo ContactRepo) *Contact {
	return &Contact{
		contactRepo: contactRepo,
	}
}

func (c *Contact) Create(createDto *contactServiceDto.Create) error {
	return c.contactRepo.Create(&contactRepositoryDto.Create{
		FirstName: createDto.FirstName,
		LastName:  createDto.LastName,
		Phones:    createDto.Phones,
	})
}
