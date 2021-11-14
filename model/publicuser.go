//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"time"
)

type PublicUser struct {
	UserID   int64      `db:"user_id,pk" json:"userId" form:"userId" query:"userId" validate:"required"`
	Fullname string     `db:"fullname,use_zero" json:"fullname" form:"fullname" query:"fullname" validate:"required,lte=80"`
	NoHp     string     `db:"no_hp,use_zero" json:"noHp" form:"noHp" query:"noHp" validate:"required,lte=20"`
	Email    string     `db:"email,use_zero" json:"email" form:"email" query:"email" validate:"required,lte=200"`
	Username string     `db:"username,use_zero" json:"username" form:"username" query:"username" validate:"required,lte=20"`
	Passwd   string     `db:"passwd,use_zero" json:"passwd" form:"passwd" query:"passwd" validate:"required,lte=200"`
	IsActive bool       `db:"is_active,use_zero" json:"isActive" form:"isActive" query:"isActive" validate:"required"`
	CreateBy int64      `db:"create_by,use_zero" json:"createBy" form:"createBy" query:"createBy" validate:"required"`
	CreateDt *time.Time `db:"create_dt,use_zero" json:"createDt" form:"createDt" query:"createDt" validate:"required"`
	UpdateBy int64      `db:"update_by,use_zero" json:"updateBy" form:"updateBy" query:"updateBy" validate:"required"`
	UpdateDt *time.Time `db:"update_dt,use_zero" json:"updateDt" form:"updateDt" query:"updateDt" validate:"required"`
	DeleteBy int64      `db:"delete_by" json:"deleteBy" form:"deleteBy" query:"deleteBy"`
	DeleteDt *time.Time `db:"delete_dt" json:"deleteDt" form:"deleteDt" query:"deleteDt"`
}
