package contact

import "github.com/labstack/echo/v4"

func (c *Handler) SetRouter(api *echo.Group) {
	api.POST("/contacts", c.Create)
	api.GET("/contacts", c.Filter)
}
