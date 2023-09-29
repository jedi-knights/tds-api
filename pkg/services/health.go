package services

import "github.com/jedi-knights/tds-api/pkg/models"

type HeathServicer interface {
	GetHealthCheck() models.HealthCheckResponse
}

type HeathService struct{}

func (s HeathService) GetHealthCheck() models.HealthCheckResponse {
	return models.HealthCheckResponse{
		Message: "Service is healthy!",
	}
}
