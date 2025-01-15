package dto

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseID struct {
	ID any `json:"id"`
}

func NewResponseID(id any) *ResponseID {
	return &ResponseID{
		ID: id,
	}
}
