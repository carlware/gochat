package rest

import (
	"context"
	"net/http"

	"github.com/carlware/gochat/accounts/cases"
	"github.com/carlware/gochat/accounts/models"
	"github.com/labstack/echo/v4"
)

func ListProfiles(c echo.Context) error {
	profiles := cases.List(context.TODO())

	return c.JSON(http.StatusOK, profiles)
}

func Login(c echo.Context) error {
	req := &cases.LoginRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Code:    "login",
			Message: err.Error(),
		})
	}

	jwt, err := cases.Login(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Code:    "login",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, jwt)
}
