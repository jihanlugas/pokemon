package db

import (
	"strconv"
	"strings"
	"time"
)

type QueryBuilder struct {
	condition    strings.Builder
	params       []interface{}
	cnt          uint8
	useOrder     bool
	hasCondition bool
}

func (q *QueryBuilder) checkPrefix() {
	if q.hasCondition {
		q.condition.WriteString(" AND ")
	} else {
		q.condition.WriteString(" WHERE ")
	}
}

func (q *QueryBuilder) Enum(field, op string, val string) *QueryBuilder {
	q.checkPrefix()
	q.hasCondition = true
	q.condition.WriteString(field)
	q.condition.WriteByte(' ')
	q.condition.WriteString(op)
	q.condition.WriteString(" $")
	q.cnt = q.cnt + 1
	q.condition.WriteString(strconv.Itoa(int(q.cnt)))
	q.params = append(q.params, val)

	return q
}

func (q *QueryBuilder) StringEq(field, value string) *QueryBuilder {
	value = strings.TrimSpace(value)
	if value != "" {
		q.checkPrefix()
		q.hasCondition = true
		q.condition.WriteString(field)
		q.condition.WriteString(" = $")
		q.cnt = q.cnt + 1
		q.condition.WriteString(strconv.Itoa(int(q.cnt)))
		q.params = append(q.params, value)
	}
	return q
}

func (q *QueryBuilder) StringLike(field, value string) *QueryBuilder {
	value = strings.TrimSpace(value)
	if value != "" {
		q.checkPrefix()
		q.hasCondition = true
		q.condition.WriteString("lower(")
		q.condition.WriteString(field)
		q.condition.WriteString(") LIKE $")
		q.cnt = q.cnt + 1
		q.condition.WriteString(strconv.Itoa(int(q.cnt)))
		q.params = append(q.params, "%"+strings.ToLower(value)+"%")
	}
	return q
}

func (q *QueryBuilder) IsNull(field string) *QueryBuilder {
	q.checkPrefix()
	q.hasCondition = true
	q.condition.WriteString(field)
	q.condition.WriteString(" IS NULL")
	return q
}

func (q *QueryBuilder) IsNotNull(field string) *QueryBuilder {
	q.checkPrefix()
	q.hasCondition = true
	q.condition.WriteString(field)
	q.condition.WriteString(" IS NOT NULL")
	return q
}

// Int64 dipakai untuk Int juga. Jika val = 0, tetap ada pecarian
// Gunakan Int64WithStr jika ada kondisi string kosong ( tidak mau ada kondisi pencarian )
// Gunakan Int64Skip Skip value 0
func (q *QueryBuilder) Int64(field, op string, val int64) *QueryBuilder {
	q.checkPrefix()
	q.hasCondition = true
	q.condition.WriteString(field)
	q.condition.WriteByte(' ')
	q.condition.WriteString(op)
	q.condition.WriteString(" $")
	q.cnt = q.cnt + 1
	q.condition.WriteString(strconv.Itoa(int(q.cnt)))
	q.params = append(q.params, val)

	return q
}

// Int64 dipakai untuk Int juga. val = 0 skip
func (q *QueryBuilder) Int64Skip(field, op string, val int64) *QueryBuilder {
	if val != 0 {
		q.checkPrefix()
		q.hasCondition = true
		q.condition.WriteString(field)
		q.condition.WriteByte(' ')
		q.condition.WriteString(op)
		q.condition.WriteString(" $")
		q.cnt = q.cnt + 1
		q.condition.WriteString(strconv.Itoa(int(q.cnt)))
		q.params = append(q.params, val)
	}

	return q
}

// Int64WithStr dipakai untuk cari yang value nya string. jika string kosong maka skip.
func (q *QueryBuilder) Int64WithStr(field, op, val string) *QueryBuilder {
	var iVal *int64
	val = strings.TrimSpace(val)

	if v, err := strconv.ParseInt(val, 10, 64); err == nil {
		iVal = &v
	}

	if iVal != nil {
		q.checkPrefix()
		q.hasCondition = true
		q.condition.WriteString(field)
		q.condition.WriteByte(' ')
		q.condition.WriteString(op)
		q.condition.WriteString(" $")
		q.cnt = q.cnt + 1
		q.condition.WriteString(strconv.Itoa(int(q.cnt)))
		q.params = append(q.params, *iVal)
	}

	return q
}

// Int64 dipakai untuk Int juga. Jika val = 0, tetap ada pecarian. Gunakan Int64WithStr jika ada kondisi string kosong ( tidak mau ada kondisi pencarian )
func (q *QueryBuilder) Timestamptz(field, op string, val *time.Time) *QueryBuilder {
	q.checkPrefix()
	q.hasCondition = true
	q.condition.WriteString(field)
	q.condition.WriteByte(' ')
	q.condition.WriteString(op)
	q.condition.WriteString(" $")
	q.cnt = q.cnt + 1
	q.condition.WriteString(strconv.Itoa(int(q.cnt)))
	q.params = append(q.params, val)

	return q
}

func (q *QueryBuilder) Bool(field string, val bool) *QueryBuilder {
	q.checkPrefix()
	q.hasCondition = true
	q.condition.WriteString(field)
	q.condition.WriteString(" = $")
	q.cnt = q.cnt + 1
	q.condition.WriteString(strconv.Itoa(int(q.cnt)))
	q.params = append(q.params, val)

	return q
}

func (q *QueryBuilder) BoolWithStr(field, val string) *QueryBuilder {
	val = strings.TrimSpace(val)
	switch val {
	case "1":
		q.Bool(field, true)
	case "0":
		q.Bool(field, false)
	}

	return q
}

func (q *QueryBuilder) Order(field, direction string) *QueryBuilder {
	if !q.useOrder {
		q.condition.WriteString(" ORDER BY ")
		q.useOrder = true
	} else {
		q.condition.WriteString(", ")
	}
	q.condition.WriteString(field)
	q.condition.WriteByte(' ')
	q.condition.WriteString(direction)
	return q
}

func (q *QueryBuilder) OffsetLimit(offset, limit int) *QueryBuilder {
	if offset > 0 {
		q.condition.WriteString(" OFFSET ")
		q.condition.WriteString(strconv.Itoa(offset))
	}
	if limit > 0 {
		q.condition.WriteString(" LIMIT ")
		q.condition.WriteString(strconv.Itoa(limit))
	}

	return q
}

func (q *QueryBuilder) Params() []interface{} {
	return q.params
}

func (q *QueryBuilder) Build() string {
	return q.condition.String()
}

type QueryComposer struct {
	sql QueryBuilder
}

func Query(s string) *QueryComposer {
	var q QueryComposer
	q.sql.condition.Grow((len(s) + 80))
	q.sql.condition.WriteString(s)
	return &q
}

func (q *QueryComposer) Where() *QueryBuilder {
	return &q.sql
}

func (q *QueryComposer) Build() string {
	return q.sql.Build()
}

func (q *QueryComposer) Params() []interface{} {
	return q.sql.params
}
