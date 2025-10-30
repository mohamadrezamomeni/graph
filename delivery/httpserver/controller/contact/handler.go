package contact

import (
	serviceDto "github.com/mohamadrezamomeni/graph/dto/service/contact"
	contactValidator "github.com/mohamadrezamomeni/graph/validator/contact"
)

type Handler struct {
	contactSvc       ContactService
	contactValidator *contactValidator.Validator
}

type ContactService interface {
	Create(*serviceDto.Create) error
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
