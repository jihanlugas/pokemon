package controller

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"gopokemon/config"
	"gopokemon/cryption"
	"gopokemon/db"
	"gopokemon/model"
	"gopokemon/response"
	"gopokemon/utils"
	"net/http"
	"time"
)

type User struct{}

func UserComposer() User {
	return User{}
}

type signinReq struct {
	Username string `db:"username,use_zero" json:"username" form:"username" query:"username" validate:"required,lte=20"`
	Passwd   string `db:"passwd,use_zero" json:"passwd" form:"passwd" query:"passwd" validate:"required,lte=200"`
}

type signupReq struct {
	Fullname     string                `json:"fullname" form:"fullname" validate:"required,lte=80"`
	Email        string                `json:"email" form:"email" validate:"required,lte=200,email,notexists=email email"`
	NoHp         string                `json:"noHp" form:"noHp" validate:"required,lte=20,notexists=no_hp noHp"`
	Username     string                `json:"username" form:"username" validate:"required,lte=20,lowercase,notexists=username username"`
	Passwd       string                `json:"passwd" form:"passwd" validate:"required,lte=200"`
}

type meRes struct {
	UserID   int64  `json:"userId"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// @Summary Sign in a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param req body signinReq true "json req body"
// @Success 200 {object} response.SuccessResponse "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /sign-in [post]
func (h User) SignIn(c echo.Context) error {
	var err error
	var ok bool

	req := new(signinReq)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(response.ResponseValidationFailed, response.ValidationError(err)).SendJSON(c)
	}

	conn, ctx, closeConn := db.GetConnection()
	defer closeConn()

	var user model.PublicUser
	user.Username = req.Username
	err = user.GetByUsername(ctx, conn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return response.Error(response.ResponseInvalidUsernameOrPassword, response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	if !user.IsActive {
		return response.Error("User not active", response.Payload{}).SendJSON(c)
	}

	if ok = cryption.CheckPasswordHash(req.Passwd, user.Passwd); !ok {
		return response.Error(response.ResponseInvalidUsernameOrPassword, response.Payload{}).SendJSON(c)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		errorInternal(c, err)
	}
	defer db.DeferHandleTransaction(ctx, tx)

	now := time.Now()

	if err = tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		errorInternal(c, err)
	}

	expiredAt := time.Now().Add(time.Hour * 12)
	token, err := getLoginToken(user.UserID, expiredAt)
	if err != nil {
		return response.Error(response.ResponseValidationFailed, response.ListErrorComposer().
			StackError("passwd", "Nama pengguna atau kata sandi tidak valid").Build()).SendJSON(c)
	}
	maxAge := expiredAt.Sub(now).Seconds()

	cookie := generateCookie(config.CookieAuthName, string(token), expiredAt, int(maxAge))
	c.SetCookie(cookie)

	return response.Success("Signin Success", response.Payload{}).SendJSON(c)
}



// @Summary Sign up a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param req body signupReq true "json req body"
// @Success 200 {object} response.SuccessResponse{payload=model.UserRes} "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /sign-up [post]
func (h User) SignUp(c echo.Context) error {
	var err error
	req := new(signupReq)

	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(response.ResponseValidationFailed, response.ValidationError(err)).SendJSON(c)
	}

	conn, ctx, closeConn := db.GetConnection()
	defer closeConn()

	newPasswd, err := cryption.HashPassword(req.Passwd)
	if err != nil {
		return response.Error("Wrong password", response.ListErrorComposer().
			StackError("passwd", "Salah format password").
			Build()).SendJSON(c)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		errorInternal(c, err)
	}
	defer db.DeferHandleTransaction(ctx, tx)

	var user model.PublicUser
	user.Fullname = req.Fullname
	user.Email = req.Email
	user.Username = req.Username
	user.NoHp = utils.FormatPhoneTo62(req.NoHp)
	user.Passwd = newPasswd
	user.IsActive = true
	user.CreateBy = 0
	user.UpdateBy = 0
	err = user.Insert(ctx, tx)
	if err != nil {
		errorInternal(c, err)
	}

	if err = tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		errorInternal(c, err)
	}

	return response.Success("Success Signup", user.UserRes()).SendJSON(c)
}



// @Summary Sign out a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /sign-out [get]
func (h User) SignOut(c echo.Context) error {
	var err error
	var user model.PublicUser
	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		return err
	}

	conn, ctx, closeConn := db.GetConnection()
	defer closeConn()

	tx, err := conn.Begin(ctx)
	if err != nil {
		errorInternal(c, err)
	}
	defer db.DeferHandleTransaction(ctx, tx)

	user.UserID = loginUser.UserID
	err = user.GetById(ctx, conn)
	if err != nil {
		errorInternal(c, err)
	}

	if err = tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		errorInternal(c, err)
	}

	expiredAt := time.Now().Add(-100 * time.Hour)
	now := time.Now()
	maxAge := expiredAt.Sub(now).Seconds()

	cookie := generateCookie(config.CookieAuthName, "", expiredAt, int(maxAge))
	c.SetCookie(cookie)

	return response.Success("Logout success!", response.Payload{}).SendJSON(c)
}


// @Tags User
// @Summary To do get current login user
// @Accept json
// @Produce json
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Success 200 {object} response.SuccessResponse{payload=meRes} "json with success = true"
// @Failure 400 {object} response.ErrorResponse "json with error = true"
// @Router /user/me [get]
func (h User) Me(c echo.Context) error {
	var err error
	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		return err
	}

	conn, ctx, closeConn := db.GetConnection()
	defer closeConn()

	var user model.PublicUser
	user.UserID = loginUser.UserID
	err = user.GetById(ctx, conn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return response.ErrorForce(response.ResponseDataNotFound, response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	res := meRes{
		UserID:   user.UserID,
		Fullname: user.Fullname,
		Email:    user.Email,
		Username: user.Username,
	}

	return response.Success("Success", res).SendJSON(c)
}

func generateCookie(Name, Token string, ExpiredAt time.Time, MaxAge int) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = Name
	cookie.Value = Token
	cookie.Path = "/"
	cookie.Expires = ExpiredAt
	cookie.MaxAge = MaxAge

	if config.Environment == config.PRODUCTION {
		cookie.SameSite = http.SameSiteNoneMode
		cookie.HttpOnly = true
		cookie.Secure = true
	} else {
		cookie.SameSite = http.SameSiteNoneMode
		cookie.HttpOnly = true
		cookie.Secure = true
	}

	return cookie
}