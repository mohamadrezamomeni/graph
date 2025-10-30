package contact

import (
	"net/http"

	"github.com/labstack/echo/v4"
	controllerDto "github.com/mohamadrezamomeni/graph/dto/controller/contact"
	serviceDto "github.com/mohamadrezamomeni/graph/dto/service/contact"
	httpErr "github.com/mohamadrezamomeni/graph/pkg/http_error"
)

func (h *Handler) Create(c echo.Context) error {
	var req controllerDto.Create
	if err := c.Bind(&req); err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	if err := h.contactValidator.ValidateCreating(req); err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.contactSvc.Create(&serviceDto.Create{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phones:    req.Phones,
	})
	if err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	return c.NoContent(http.StatusNoContent)
}
