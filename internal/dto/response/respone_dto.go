package response

type ResponseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error"`
}

func FormaterResponseDTO(code int, message string, data interface{}, err interface{}) *ResponseModel {
	return &ResponseModel{
		Code:    code,
		Message: message,
		Data:    data,
		Error:   err,
	}
}

type PaginationResponseDTO struct {
	TotalItems  int         `json:"total_items"`
	CurrentPage int         `json:"current_page"`
	Items       interface{} `json:"items"`
}

func FormaterPaginationResponseDTO(totalItems int, currentPage int, items interface{}) *PaginationResponseDTO {
	return &PaginationResponseDTO{
		TotalItems:  totalItems,
		CurrentPage: currentPage,
		Items:       items,
	}
}
