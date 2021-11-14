//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopokemon/db"
	"strings"
	"time"
)

type UserRes struct {
	UserID   int64  `json:"userId"`
	Fullname string `json:"fullname"`
	NoHp     string `json:"noHp"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
}

func GetUserQuery() *db.QueryComposer {
	return db.Query(`SELECT user_id, fullname, email, username, no_hp, passwd, is_active, create_by, create_dt, update_by, update_dt FROM public.user`)
}

func (p *PublicUser) GetById(ctx context.Context, conn *pgxpool.Conn) error {
	var err error

	sql := GetUserQuery().
		Where().
		Int64(`user_id`, "=", int64(p.UserID)).
		IsNull(`delete_dt`).
		OffsetLimit(0, 1)
	err = pgxscan.Get(ctx, conn, p, sql.Build(), sql.Params()...)

	return err
}

func (p *PublicUser) GetByUsername(ctx context.Context, conn *pgxpool.Conn) error {
	var err error
	p.Username = strings.ToLower(p.Username)

	sql := GetUserQuery().
		Where().
		StringEq(`username`, p.Username).
		IsNull(`delete_dt`).
		OffsetLimit(0, 1)
	err = pgxscan.Get(ctx, conn, p, sql.Build(), sql.Params()...)

	return err
}

func (p *PublicUser) Insert(ctx context.Context, tx pgx.Tx) error {
	var err error

	now := time.Now()
	p.Username = strings.ToLower(p.Username)
	p.CreateDt = &now
	p.UpdateDt = &now
	err = tx.QueryRow(ctx, `INSERT INTO public.user (fullname, email, username, no_hp, passwd, is_active, create_by, create_dt, update_by, update_dt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING user_id;`,
		p.Fullname,
		p.Email,
		p.Username,
		p.NoHp,
		p.Passwd,
		p.IsActive,
		p.CreateBy,
		p.CreateDt,
		p.UpdateBy,
		p.UpdateDt,
	).Scan(&p.UserID)
	return err
}

func (p *PublicUser) UserRes() UserRes {
	var res UserRes

	res.UserID = p.UserID
	res.Fullname = p.Fullname
	res.NoHp = p.NoHp
	res.Email = p.Email
	res.IsActive = p.IsActive

	return res
}

func ToUsersRes(users []PublicUser) []UserRes {
	res := make([]UserRes, 0)

	for _, data := range users {
		res = append(res, data.UserRes())
	}

	return res
}
