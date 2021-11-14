package controller

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"gopokemon/db"
	"gopokemon/model"
	"gopokemon/response"
)

type Userpokemon struct{}

func UserpokemonComposer() Userpokemon {
	return Userpokemon{}
}

type mypokemonRes struct {
	UserpokemonID int64      `json:"userpokemonId"`
	UserID        int64      `json:"userId"`
	Pokemon       string     `json:"pokemon"`
	Nickname      string     `json:"nickname"`
}

type updateUserpokemonReq struct {
	UserpokemonID int64      `json:"userpokemonId" validate:"required"`
	Nickname      string     `json:"nickname" validate:"required"`
}

// @Tags User Pokemon
// @Summary page user pokemon
// @Accept json
// @Produce json
// @in header
// @Success 200 {object} response.SuccessResponse{payload=[]model.UserpokemonRes} "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /user-pokemon/my-pokemon [get]
func (h Userpokemon) MyPokemon(c echo.Context) error {
	var err error
	var listUserpokemon []model.PublicUserpokemon

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	conn, ctx, closeConn := db.GetConnection()
	defer closeConn()

	q := model.GetUserPokemonQuery().Where().
		Int64("user_id", "=", loginUser.UserID)
	listUserpokemon, err = model.GetUserpokemonWhere(ctx, conn, q)
	if err != nil {
		errorInternal(c, err)
	}

	res := model.ToUserpokemonRes(listUserpokemon)

	return response.Success(response.ResponseSuccess, res).SendJSON(c)
}

// @Tags User Pokemon
// @Summary update user pokemon
// @Accept json
// @Produce json
// @in header
// @Param req body updateUserpokemonReq true "Req Param"
// @Success 200 {object} response.SuccessResponse{payload=model.UserpokemonRes} "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /user-pokemon/my-pokemon [post]
func (h Userpokemon) Update(c echo.Context) error {
	var err error
	var userpokemon model.PublicUserpokemon

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	req := new(updateUserpokemonReq)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(response.ResponseValidationFailed, response.ValidationError(err)).SendJSON(c)
	}

	conn, ctx, closeConn := db.GetConnection()
	defer closeConn()

	userpokemon.UserpokemonID = req.UserpokemonID
	err = userpokemon.GetById(ctx, conn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return response.Error(response.ResponseDataNotFound, response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		errorInternal(c, err)
	}
	defer db.DeferHandleTransaction(ctx, tx)

	userpokemon.Nickname = req.Nickname
	userpokemon.UpdateBy = loginUser.UserID
	err = userpokemon.Update(ctx, tx)
	if err != nil {
		errorInternal(c, err)
	}

	if err = tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		errorInternal(c, err)
	}

	return response.Success(response.ResponseSuccess, userpokemon.UserpokemonRes()).SendJSON(c)

}

