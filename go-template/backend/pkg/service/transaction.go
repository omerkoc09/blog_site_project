package service

import (
	"github.com/uptrace/bun"
)

type IServiceTransaction[S any] interface {
	SetTx(tx bun.Tx) S // return service
}

// örnek kullanım:
// iligli servisin interface'ine IServiceTransaction eklenir

//type IXService interface {
//	service.IBaseService[model.X]
//	service.IServiceTransaction[IXService]
//}
//
// ilgili service struct'ına bunun implementasyonu eklenir
// parametre olarak gelen tx bunun db'si olarak set edilir

//func (s XService) SetTx(tx bun.Tx) IXService {
//	s.DB = tx
//	return s
//}
