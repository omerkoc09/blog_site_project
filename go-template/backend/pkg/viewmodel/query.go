package viewmodel

type QueryColumnType int

const (
	QueryColumnTypeNumber            QueryColumnType = 1
	QueryColumnTypeString            QueryColumnType = 10
	QueryColumnTypeBool              QueryColumnType = 20
	QueryColumnTypeJsonbNumber       QueryColumnType = 30
	QueryColumnTypeJsonbString       QueryColumnType = 40
	QueryColumnTypeAssociationNumber QueryColumnType = 50
	QueryColumnTypeAssociationString QueryColumnType = 60
	QueryColumnTypeInInt64           QueryColumnType = 70
	QueryColumnTypeInString          QueryColumnType = 80
	QueryColumnTypeNull              QueryColumnType = 90
	QueryColumnTypeNotNull           QueryColumnType = 91
	QueryColumnTypeDateBuyuk         QueryColumnType = 92
)

type QueryModel struct {
	Columns     []string          `query:"columns"`
	ColumnTypes []QueryColumnType `query:"column_types"`
	Query       []string          `query:"query"`
}

func (q *QueryModel) Append(query, column string, columnType QueryColumnType) {
	q.Query = append(q.Query, query)
	q.Columns = append(q.Columns, column)
	q.ColumnTypes = append(q.ColumnTypes, columnType)
}
