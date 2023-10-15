package model

type GenericResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data,omitempty"`
}
