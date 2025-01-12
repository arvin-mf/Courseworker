package dto

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseID struct {
	ID string `json:"id"`
}

func NewResponseID(id string) *ResponseID {
	return &ResponseID{
		ID: id,
	}
}
