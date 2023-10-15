package model

type Status struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Status     string `json:"status"`
	SwaggerURL string `json:"swaggerURL"`
}
