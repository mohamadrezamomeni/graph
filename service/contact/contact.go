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
	Create(*contactRepositoryDto.Create) (string, error)
	Filter(*contactRepositoryDto.Filter) ([]*entity.Contact, error)
	Update(string, *contactRepositoryDto.Update) error
}

func New(contactRepo ContactRepo) *Contact {
	return &Contact{
		contactRepo: contactRepo,
	}
}

func (c *Contact) Create(createDto *contactServiceDto.Create) error {
	_, err := c.contactRepo.Create(&contactRepositoryDto.Create{
		FirstName: createDto.FirstName,
		LastName:  createDto.LastName,
		Phones:    createDto.Phones,
	})
	return err
}

func (c *Contact) Filter(filterDto *contactServiceDto.Filter) ([]*entity.Contact, error) {
	return c.contactRepo.Filter(&contactRepositoryDto.Filter{
		FirstNames: filterDto.FirstNames,
		LastNames:  filterDto.LastNames,
		Phones:     filterDto.Phones,
	})
}

func (c *Contact) Update(id string, updateDto *contactServiceDto.Update) error {
	return c.contactRepo.Update(id, &contactRepositoryDto.Update{
		FirstName: updateDto.FirstName,
		LastName:  updateDto.LastName,
		Phones:    updateDto.Phones,
	})
}
