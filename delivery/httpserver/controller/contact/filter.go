package contact

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	controllerDto "github.com/mohamadrezamomeni/graph/dto/controller/contact"
	serviceDto "github.com/mohamadrezamomeni/graph/dto/service/contact"
	httpErr "github.com/mohamadrezamomeni/graph/pkg/http_error"
	serializer "github.com/mohamadrezamomeni/graph/serializer/contact"
)

func (h *Handler) Filter(c echo.Context) error {
	var req controllerDto.Filter
	if err := c.Bind(&req); err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	filterServiceDto := new(serviceDto.Filter)

	if len(req.FirstNames) > 0 {
		filterServiceDto.FirstNames = strings.Split(req.FirstNames, ",")
	}
	if len(req.Phones) > 0 {
		filterServiceDto.Phones = strings.Split(req.Phones, ",")
	}
	if len(req.LastNames) > 0 {
		filterServiceDto.Phones = strings.Split(req.LastNames, ",")
	}
	contacts, err := h.contactSvc.Filter(filterServiceDto)
	if err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, msg)
	}
	concatSerializer := &serializer.FilterConcats{
		Items: make([]*serializer.Contact, 0),
	}

	for _, concat := range contacts {
		concatSerializer.Items = append(concatSerializer.Items, &serializer.Contact{
			ID:        concat.ID,
			FirstName: concat.FirstName,
			LastName:  concat.LastName,
			Phones:    concat.Phones,
		})
	}

	return c.JSON(http.StatusOK, concatSerializer)
}
