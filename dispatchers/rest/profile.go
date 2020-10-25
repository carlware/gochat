package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carlware/gochat/accounts/cases"
	"github.com/labstack/echo/v4"
)

type StandardError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ListProfiles(c echo.Context) error {
	profiles := cases.List(context.TODO())

	return c.JSON(http.StatusOK, profiles)
}

func Login(c echo.Context) error {
	req := &cases.LoginRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &StandardError{
			Code:    "login",
			Message: err.Error(),
		})
	}

	jwt, err := cases.Login(req)
	fmt.Println(jwt, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &StandardError{
			Code:    "login",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, jwt)
}
