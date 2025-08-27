package models

type GeneralResponse struct {
	ResponseStatus      string      `json:"response_status"`
	ResponseDescription string      `json:"response_description"`
	ResponseData        interface{} `json:"response_data,omitempty"`
	StatusCode          int         `json:"status_code"`
}
