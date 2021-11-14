package router

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gopokemon/config"
	"gopokemon/constant"
	"gopokemon/controller"
	"gopokemon/cryption"
	"gopokemon/response"
	"gopokemon/swg"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type host struct {
	echo *echo.Echo
}

func Init() *echo.Echo {
	var err error
	var listenPort int
	hosts := make(map[string]*host)
	if listenPort, err = strconv.Atoi(config.ListenTo.Port); err != nil {
		panic(err)
	}

	web := websiteRouter()
	webDomain := config.WebDomainName
	hosts[webDomain] = &host{web}
	hosts[fmt.Sprintf("%s:%v", webDomain, listenPort)] = hosts[webDomain]

	checkToken := checkTokenMiddleware()

	//userController := controller.UserComposer()
	//mskanjiController := controller.MskanjiComposer()
	//kanjiController := controller.KanjiComposer()
	//mskanjiexampleController := controller.MskanjiexampleComposer()
	//// propertyController := controller.PropertyComposer()
	//// productController := controller.ProductComposer()
	//// itemController := controller.ItemComposer()


	userController := controller.UserComposer()
	pokemonController := controller.PokemonComposer()
	userpokemonController := controller.UserpokemonComposer()


	//if config.Environment != config.PRODUCTION {
		web.GET("/swg/*", echoSwagger.WrapHandler)
	//}

	webPokemon := web.Group("/pokemon", checkToken)
	webPokemon.GET("", pokemonController.Page)
	webPokemon.GET("/:pokemon", pokemonController.Get)
	webPokemon.POST("/catch", pokemonController.Catch)

	web.GET("/", controller.Ping)
	web.POST("/sign-up", userController.SignUp)
	web.POST("/sign-in", userController.SignIn)
	web.GET("/sign-out", userController.SignOut, checkToken)

	webUser := web.Group("/user", checkToken)
	webUser.GET("/me", userController.Me)

	webUserpokemon := web.Group("/user-pokemon", checkToken)
	webUserpokemon.GET("/my-pokemon", userpokemonController.MyPokemon)
	webUserpokemon.POST("/my-pokemon", userpokemonController.Update)

	//
	//webUser := web.Group("/user", checkToken)
	//webUser.GET("/me", userController.Me)
	//
	//webMskanji := web.Group("/mskanji", checkToken)
	//webMskanji.GET("/fetch", mskanjiController.Fetch)
	//webMskanji.GET("/search", mskanjiController.Search)
	//webMskanji.GET("/addtokanji/:mskanji_id", mskanjiController.AddToKanji)
	//webMskanji.POST("/update", mskanjiController.Update)
	//
	//webKanji := web.Group("/kanji", checkToken)
	//webKanji.POST("/update", kanjiController.Update)
	//webKanji.GET("/list", kanjiController.List)
	//webKanji.GET("/:kanji", kanjiController.GetByKanji)
	//webKanji.GET("/id/:kanji_id", kanjiController.Get)
	//
	//webMskanjiexample := web.Group("/mskanji-example", checkToken)
	//webMskanjiexample.GET("/:mskanjiexampleId", mskanjiexampleController.Get)
	//webMskanjiexample.POST("/list", mskanjiexampleController.List)
	//webMskanjiexample.POST("/create", mskanjiexampleController.Create)

	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		hst := hosts[req.Host]

		//if config.Environment != config.PRODUCTION {
		swg.SwaggerInfo.Title = "Website API"
		swg.SwaggerInfo.BasePath = "/"
		//}

		if hst == nil {
			err = echo.ErrNotFound
		} else {
			hst.echo.ServeHTTP(res, req)
		}

		return
	})

	return e
}

func checkTokenMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(config.CookieAuthName)
			if err != nil {
				return response.ErrorForce("Akses ditolak!", response.Payload{}).SendJSON(c)
			}

			token := cookie.Value
			tokenPayload, err := cryption.DecryptAES64([]byte(token))
			if err != nil {
				return response.ErrorForce("Akses telah kadaluarsa", response.Payload{}).SendJSON(c)
			}

			if len(tokenPayload) == constant.TokenPayloadLen {
				expiredUnix := binary.BigEndian.Uint64(tokenPayload)
				expiredAt := time.Unix(int64(expiredUnix), 0)
				now := time.Now()
				if now.After(expiredAt) {
					return response.ErrorForce("Akses telah kadaluarsa!", response.Payload{}).SendJSON(c)
				} else {
					usrLogin := controller.UserLogin{
						UserID:      int64(binary.BigEndian.Uint64(tokenPayload[8:])),
					}
					c.Set(constant.TokenUserContext, usrLogin)
					return next(c)
				}
			} else {
				return response.ErrorForce("Akses telah kadaluarsa!", response.Payload{}).SendJSON(c)
			}
		}
	}
}

func httpErrorHandler(err error, c echo.Context) {
	var errorResponse *response.ErrorResponse
	code := http.StatusInternalServerError

	switch e := err.(type) {
	case *echo.HTTPError:
		// Handle pada saat URL yang di request tidak ada. atau ada kesalahan server.
		code = e.Code
		errorResponse = &response.ErrorResponse{
			IsError: true,
			Message: strconv.Itoa(code) + " code. " + fmt.Sprintf("%v", e.Message),
			Payload: map[string]interface{}{},
		}
	case *response.ErrorResponse:
		errorResponse = e
	default:
		// Handle error dari panic
		if config.Environment != config.PRODUCTION {
			errorResponse = &response.ErrorResponse{
				IsError: true,
				Message: err.Error(),
				Payload: map[string]interface{}{},
			}
		} else {
			code = http.StatusInternalServerError
			errorResponse = &response.ErrorResponse{
				IsError: true,
				Message: "Internal server error",
				Payload: map[string]interface{}{},
			}
		}
	}

	js, err := json.Marshal(errorResponse)
	if err == nil {
		c.String(code, string(js))
	} else {
		b := []byte("{error: true, message: \"unresolved error\"}")
		c.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, b)
	}
}
