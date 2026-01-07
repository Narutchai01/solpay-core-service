package response

type ResponseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error"`
}

func NewResponseModel(code int, message string, data interface{}, err interface{}) *ResponseModel {
	return &ResponseModel{
		Code:    code,
		Message: message,
		Data:    data,
		Error:   err,
	}
}
