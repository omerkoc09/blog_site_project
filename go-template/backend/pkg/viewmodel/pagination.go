package viewmodel

type SortColumnType int

const (
	SortColumnTypeNormal      SortColumnType = 1
	SortColumnTypeJsonb       SortColumnType = 10
	SortColumnTypeAssociation SortColumnType = 20
)

type SortOrder string

const (
	SortOrderASC  SortOrder = "ASC"
	SortOrderDESC SortOrder = "DESC"
)

type PaginationModel struct {
	Page            int              `query:"page"`
	PerPage         int              `query:"per_page"`
	Offset          int              `query:"offset"`
	SortColumns     []string         `query:"sort_columns"`
	SortColumnTypes []SortColumnType `query:"sort_column_types"`
	SortOrders      []SortOrder      `query:"sort_orders"`
}
