package accounts

import (
	"github.com/carlware/gochat/common/auth"
	"github.com/carlware/gochat/dispatchers/accounts/rest"
	"github.com/labstack/echo/v4"
)

func RunMicroAccounts(e *echo.Echo) {

	e.GET("/profiles", rest.ListProfiles, auth.IsLoggedIn)
	e.POST("/login", rest.Login)
}
