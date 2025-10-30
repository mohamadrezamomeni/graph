package contact

import (
	contactRepositoryDto "github.com/mohamadrezamomeni/graph/dto/repository/contact"
	contactServiceDto "github.com/mohamadrezamomeni/graph/dto/service/contact"
	"github.com/mohamadrezamomeni/graph/entity"
)

type Contact struct {
	contactRepo ContactRepo
}

type ContactRepo interface {
	Create(*contactRepositoryDto.Create) error
	Filter(*contactRepositoryDto.Filter) ([]*entity.Contact, error)
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

func (c *Contact) Filter(filterDto *contactServiceDto.Filter) ([]*entity.Contact, error) {
	return c.contactRepo.Filter(&contactRepositoryDto.Filter{
		FirstNames: filterDto.FirstNames,
		LastNames:  filterDto.LastNames,
		Phones:     filterDto.Phones,
	})
}
