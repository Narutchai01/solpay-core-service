package response

// ResponseModel is the standard API response envelope.
type ResponseModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error"`
}

// FormatResponseDTO creates a standard API response.
func FormatResponseDTO(code int, message string, data any, err any) *ResponseModel {
	return &ResponseModel{
		Code:    code,
		Message: message,
		Data:    data,
		Error:   err,
	}
}

// PaginationResponseDTO wraps a paginated list response.
type PaginationResponseDTO struct {
	TotalItems  int `json:"total_items"`
	CurrentPage int `json:"current_page"`
	Items       any `json:"items"`
}

// FormatPaginationResponseDTO creates a paginated response.
func FormatPaginationResponseDTO(totalItems int, currentPage int, items any) *PaginationResponseDTO {
	return &PaginationResponseDTO{
		TotalItems:  totalItems,
		CurrentPage: currentPage,
		Items:       items,
	}
}
