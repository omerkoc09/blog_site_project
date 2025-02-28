package viewmodel

type ResponseModel struct {
	Data         interface{} `json:"data"`
	DataCount    int         `json:"data_count"`
	ErrorMessage string      `json:"error_message"`
	ErrorCode    int         `json:"error_code"`
}

type Validation interface {
	Validate() []error
}
