package contact

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	contactControllerDto "github.com/mohamadrezamomeni/graph/dto/controller/contact"
	appError "github.com/mohamadrezamomeni/graph/pkg/error"
)

func (v *Validator) ValidateUpdating(req contactControllerDto.Update) error {
	scope := "validator.contact.validateCreating"

	err := validation.ValidateStruct(&req,
		validation.Field(&req.FirstName,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[A-Za-z]{2,}$`))),

		validation.Field(&req.LastName,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[A-Za-z]{2,}$`))),

		validation.Field(&req.Phones,
			validation.Required,
			validation.Each(
				validation.Required,
				validation.Match(regexp.MustCompile("^09[0-9]{9}$")).Error("error to validate phone"),
			),
		),
	)
	if err != nil {
		return appError.Wrap(err).Scope(scope).BadRequest().Input(req).ErrorWrite()
	}

	return nil
}
