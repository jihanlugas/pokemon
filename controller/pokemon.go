package controller

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/mtslzr/pokeapi-go"
	"gopokemon/db"
	"gopokemon/model"
	"gopokemon/response"
	"gopokemon/structs"
	"math/rand"
)

type Pokemon struct{}

func PokemonComposer() Pokemon {
	return Pokemon{}
}

type pageReq struct {
	Limit  int `json:"limit" query:"limit"`
	Offset int `json:"offset" query:"offset"`
}

type result struct {
	structs.Result
	Image string `json:"image"`
}

type pageRes struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []result    `json:"results"`
}

type catchReq struct {
	Pokemon string `json:"pokemon" validate:"required"`
}

func getPage(data structs.Resource) pageRes{
	var res pageRes

	res.Next = data.Next
	res.Count = data.Count
	res.Previous = data.Previous

	for i, datares := range data.Results {
		res.Results[i].Name = datares.Name
		res.Results[i].URL = datares.URL

		//detail, err := pokeapi.Pokemon("bulbasaur")
		//if err != nil {
		//	continue
		//}
		res.Results[i].Image = ""
	}

	return res
}

// @Tags Pokemon
// @Summary page pokemon
// @Accept json
// @Produce json
// @in header
// @Param req query pageReq true "Req Param"
// @Success 200 {object} response.SuccessResponse{payload=pageRes} "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /pokemon [get]
func (h Pokemon) Page(c echo.Context) error {
	var err error

	req := new(pageReq)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(response.ResponseValidationFailed, response.ValidationError(err)).SendJSON(c)
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	list, err := pokeapi.Resource("pokemon", req.Offset, req.Limit)
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(response.ResponseSuccess, list).SendJSON(c)
}

// @Tags Pokemon
// @Summary get pokemon
// @Accept json
// @Produce json
// @in header
// @name Authorization
// @Param pokemon path string true "Pokemon"
// @Success 200 {object} response.SuccessResponse{payload=structs.Pokemon} "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /pokemon/{pokemon} [get]
func (h Pokemon) Get(c echo.Context) error {
	var err error
	req := c.Param("pokemon")
	if err != nil {
		return response.Error("Not Found", response.Payload{}).SendJSON(c)
	}

	pokemon, err := pokeapi.Pokemon(req)
	if err != nil {
		return response.Error(response.ResponseDataNotFound, response.Payload{}).SendJSON(c)
	}

	return response.Success(response.ResponseSuccess, pokemon).SendJSON(c)
}



// @Tags Pokemon
// @Summary catch pokemon with success probability 50%
// @Accept json
// @Produce json
// @in header
// @name Authorization
// @Param req body catchReq true "Req Param"
// @Success 200 {object} response.SuccessResponse{payload=model.UserpokemonRes} "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /pokemon/catch [post]
func (h Pokemon) Catch(c echo.Context) error {
	var err error
	var userpokemon model.PublicUserpokemon

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	req := new(catchReq)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(response.ResponseValidationFailed, response.ValidationError(err)).SendJSON(c)
	}

	pokemon, err := pokeapi.Pokemon(req.Pokemon)
	if err != nil {
		return response.Error(response.ResponseDataNotFound, response.Payload{}).SendJSON(c)
	}

	conn, ctx, closeConn := db.GetConnection()
	defer closeConn()

	userpokemon.UserID = loginUser.UserID
	userpokemon.Pokemon = pokemon.Name
	err = userpokemon.GetByUserPokemon(ctx, conn)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			errorInternal(c, err)
		}
	} else {
		return response.Error("anda_sudah_punya_pokemon_ini", response.Payload{}).SendJSON(c)
	}

	// mengakap pokemon dengan probability success 50%
	if rand.Int63n(8999)%2 == 0 {
		return response.Error("anda_gagal_menangkap_pokemon_ini_coba_lagi", response.Payload{}).SendJSON(c)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		errorInternal(c, err)
	}
	defer db.DeferHandleTransaction(ctx, tx)

	err = userpokemon.Insert(ctx, tx)
	if err != nil {
		errorInternal(c, err)
	}

	if err = tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		errorInternal(c, err)
	}

	return response.Success(response.ResponseSuccess, userpokemon.UserpokemonRes()).SendJSON(c)
}




