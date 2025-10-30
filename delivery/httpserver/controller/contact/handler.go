package contact

import (
	serviceDto "github.com/mohamadrezamomeni/graph/dto/service/contact"
	"github.com/mohamadrezamomeni/graph/entity"
	contactValidator "github.com/mohamadrezamomeni/graph/validator/contact"
)

type Handler struct {
	contactSvc       ContactService
	contactValidator *contactValidator.Validator
}

type ContactService interface {
	Create(*serviceDto.Create) error
	Filter(*serviceDto.Filter) ([]*entity.Contact, error)
}

func New(
	contactSvc ContactService,
	contactValidator *contactValidator.Validator,
) *Handler {
	return &Handler{
		contactSvc:       contactSvc,
		contactValidator: contactValidator,
	}
}
