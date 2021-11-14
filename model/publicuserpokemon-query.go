//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopokemon/db"
	"time"
)

type UserpokemonRes struct {
	UserpokemonID int64      `db:"userpokemon_id,pk" json:"userpokemonId"`
	UserID        int64      `db:"user_id,use_zero" json:"userId"`
	Pokemon       string     `db:"pokemon,use_zero" json:"pokemon"`
	Nickname      string     `db:"nickname,use_zero" json:"nickname"`
}

func GetUserPokemonQuery() *db.QueryComposer {
	return db.Query(`SELECT userpokemon_id, user_id, pokemon, nickname, create_by, create_dt, update_by, update_dt FROM public.userpokemon`)
}

func (p *PublicUserpokemon) GetById(ctx context.Context, conn *pgxpool.Conn) error {
	var err error

	sql := GetUserPokemonQuery().
		Where().
		Int64(`userpokemon_id`, "=", int64(p.UserpokemonID)).
		OffsetLimit(0, 1)
	err = pgxscan.Get(ctx, conn, p, sql.Build(), sql.Params()...)

	return err
}

func (p *PublicUserpokemon) GetByUserPokemon(ctx context.Context, conn *pgxpool.Conn) error {
	var err error
	sql := GetUserPokemonQuery().
		Where().
		Int64(`user_id`, "=", p.UserID).
		StringEq(`pokemon`, p.Pokemon).
		OffsetLimit(0, 1)
	err = pgxscan.Get(ctx, conn, p, sql.Build(), sql.Params()...)

	return err
}

func (p *PublicUserpokemon) Insert(ctx context.Context, tx pgx.Tx) error {
	var err error

	now := time.Now()
	p.CreateDt = &now
	p.UpdateDt = &now
	err = tx.QueryRow(ctx, `INSERT INTO public.userpokemon (user_id, pokemon, nickname, create_by, create_dt, update_by, update_dt)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING userpokemon_id;`,
		p.UserID,
		p.Pokemon,
		p.Nickname,
		p.CreateBy,
		p.CreateDt,
		p.UpdateBy,
		p.UpdateDt,
	).Scan(&p.UserpokemonID)
	return err
}

func (p *PublicUserpokemon) Update(ctx context.Context, tx pgx.Tx) error {
	var err error

	now := time.Now()
	p.UpdateDt = &now
	_, err = tx.Exec(ctx, `UPDATE public.userpokemon SET nickname = $1
		, update_by = $2
		, update_dt = $3
		WHERE userpokemon_id = $4`,
		p.Nickname,
		p.UpdateBy,
		p.UpdateDt,
		p.UserpokemonID,
	)
	return err
}

func (p *PublicUserpokemon) UserpokemonRes() UserpokemonRes {
	var res UserpokemonRes

	res.UserpokemonID = p.UserpokemonID
	res.UserID = p.UserID
	res.Pokemon = p.Pokemon
	res.Nickname = p.Nickname

	return res
}

func GetUserpokemonWhere(ctx context.Context, conn *pgxpool.Conn, q *db.QueryBuilder) ([]PublicUserpokemon, error) {
	var err error
	var data []PublicUserpokemon

	err = pgxscan.Select(ctx, conn, &data, q.Build(), q.Params()...)
	if err != nil {
		return data, err
	}
	if len(data) == 0 {
		data = make([]PublicUserpokemon, 0)
	}

	return data, err
}

func ToUserpokemonRes(items []PublicUserpokemon) []UserpokemonRes {
	res := make([]UserpokemonRes, 0)

	for _, data := range items {
		res = append(res, data.UserpokemonRes())
	}

	return res
}




