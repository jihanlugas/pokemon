//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"time"
)

type PublicUserpokemon struct {
	UserpokemonID int64      `db:"userpokemon_id,pk" json:"userpokemonId" form:"userpokemonId" query:"userpokemonId" validate:"required"`
	UserID        int64      `db:"user_id,use_zero" json:"userId" form:"userId" query:"userId" validate:"required"`
	Pokemon       string     `db:"pokemon,use_zero" json:"pokemon" form:"pokemon" query:"pokemon" validate:"required,lte=80"`
	Nickname      string     `db:"nickname,use_zero" json:"nickname" form:"nickname" query:"nickname" validate:"required,lte=80"`
	CreateBy      int64      `db:"create_by,use_zero" json:"createBy" form:"createBy" query:"createBy" validate:"required"`
	CreateDt      *time.Time `db:"create_dt,use_zero" json:"createDt" form:"createDt" query:"createDt" validate:"required"`
	UpdateBy      int64      `db:"update_by,use_zero" json:"updateBy" form:"updateBy" query:"updateBy" validate:"required"`
	UpdateDt      *time.Time `db:"update_dt,use_zero" json:"updateDt" form:"updateDt" query:"updateDt" validate:"required"`
}
