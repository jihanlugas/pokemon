package router

import (
	"gopokemon/controller"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func websiteRouter() *echo.Echo {
	e := echo.New()
	e.Validator = controller.Validate
	e.HTTPErrorHandler = httpErrorHandler
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll:   true, // config.Env == config.PRODUCTION
		DisablePrintStack: true, // config.Env == config.PRODUCTION
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*", "https://pos.nextkanji.com"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderCookie, echo.HeaderXRequestedWith, echo.HeaderXRealIP, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderSetCookie},
		Skipper: func(c echo.Context) bool {
			return false
		},
	}))

	return e
}
