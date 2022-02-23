package service

import "parserProject/models"

// ServiceInterface Основной интерфейс для эндпоинтс
type ServiceInterface interface {
	GetTitlesWithUrl(urls []string) ([]models.ResponseTitlesAndUrls, error)
}
