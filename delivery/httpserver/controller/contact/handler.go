package contact

import (
	contactControllerDto "github.com/mohamadrezamomeni/graph/dto/controller/contact"
	serviceDto "github.com/mohamadrezamomeni/graph/dto/service/contact"
	"github.com/mohamadrezamomeni/graph/entity"
)

type Handler struct {
	contactSvc       ContactService
	contactValidator ContactValidation
}

type ContactService interface {
	Create(*serviceDto.Create) error
	Filter(*serviceDto.Filter) ([]*entity.Contact, error)
	Update(string, *serviceDto.Update) error
}

type ContactValidation interface {
	ValidateUpdating(contactControllerDto.Update) error
	ValidateCreating(contactControllerDto.Create) error
}

func New(
	contactSvc ContactService,
	contactValidator ContactValidation,
) *Handler {
	return &Handler{
		contactSvc:       contactSvc,
		contactValidator: contactValidator,
	}
}
