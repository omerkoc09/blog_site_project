package service

import (
	"strings"

	"github.com/uptrace/bun"

	"github.com/hayrat/go-template2/backend/pkg/viewmodel"
)

func Filter(sq *bun.SelectQuery, f *viewmodel.QueryModel) *bun.SelectQuery {
	for i, column := range f.Columns {
		switch f.ColumnTypes[i] {
		case viewmodel.QueryColumnTypeNumber, viewmodel.QueryColumnTypeBool, viewmodel.QueryColumnTypeAssociationNumber:
			sq = sq.Where("? = ?", bun.Ident(column), f.Query[i])
		case viewmodel.QueryColumnTypeString, viewmodel.QueryColumnTypeAssociationString:
			columns := strings.Split(column, ";")
			sq = sq.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				for _, c := range columns {
					q = q.WhereOr("? ILIKE ?", bun.Ident(c), "%"+f.Query[i]+"%")
				}
				return q
			})
		case viewmodel.QueryColumnTypeDateBuyuk:
			sq = sq.Where("? > ?", bun.Ident(column), f.Query[i])
		case viewmodel.QueryColumnTypeInString, viewmodel.QueryColumnTypeInInt64:
			sq = sq.Where("? IN (?)", bun.Ident(column), bun.In(split(f.Query[i])))
		case viewmodel.QueryColumnTypeJsonbNumber:
			c := strings.Split(column, ".")
			sq = sq.Where("(?->>"+"'"+c[1]+"'"+")::int"+" = "+f.Query[i], bun.Ident(c[0]))
		case viewmodel.QueryColumnTypeJsonbString:
			c := strings.Split(column, ".")
			sq = sq.Where("? ->>"+"'"+c[1]+"'"+" ILIKE '%"+f.Query[i]+"%'", bun.Ident(c[0]))
		case viewmodel.QueryColumnTypeNull:
			sq = sq.Where("? IS NULL", bun.Ident(column))
		case viewmodel.QueryColumnTypeNotNull:
			sq = sq.Where("? IS NOT NULL", bun.Ident(column))
		}
	}

	return sq
}

func Paginate(sq *bun.SelectQuery, p *viewmodel.PaginationModel) *bun.SelectQuery {
	for i, column := range p.SortColumns {
		switch p.SortColumnTypes[i] {
		case viewmodel.SortColumnTypeNormal:
			sq = sq.OrderExpr("? ?", bun.Ident(column), bun.Safe(p.SortOrders[i]))
		case viewmodel.SortColumnTypeJsonb:
			c := strings.Split(column, ".")
			sq = sq.OrderExpr("? ->> ? ?", bun.Ident(c[0]), bun.Safe(c[2]), bun.Safe(p.SortOrders[i]))
		case viewmodel.SortColumnTypeAssociation:
			sq = sq.OrderExpr("? ?", bun.Ident(column), bun.Safe(p.SortOrders[i]))
		}
	}

	return sq.Offset(p.Offset).Limit(p.PerPage)
}

func split(s string) []string {
	s = strings.ReplaceAll(s, ",", ";")
	s = strings.ReplaceAll(s, " ", ";")
	return strings.Split(s, ";")
}
