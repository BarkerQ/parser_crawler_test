package service

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
	"parserProject/httpClients"
	"parserProject/models"
)

type service struct {
	logger          log.Logger
	clientInterface httpClients.HttpClientInterface
}

func NewService(logger log.Logger, clientInterface httpClients.HttpClientInterface) ServiceInterface {
	return &service{
		logger:          logger,
		clientInterface: clientInterface,
	}
}

// GetTitlesWithUrl Метод получения массива ссылок и получение данных
func (s *service) GetTitlesWithUrl(urls []string) ([]models.ResponseTitlesAndUrls, error) {
	s.logger.Log("method", "GetTittleWithUrl")

	var titlesAndUrl []models.ResponseTitlesAndUrls
	for _, url := range urls {

		// Получаем данные источника по ссылке
		getDataFromSource, err := s.clientInterface.ParseHtmlFromUrlSource(url)
		if err == nil {

			// После получения данных ищем заголовок
			parse, err := getTitleFromSource(getDataFromSource, url)
			if err == nil {

				// Если все ок и заголовок нашелся, добавляем в конечный массив
				titlesAndUrl = append(titlesAndUrl, parse)
			}
		}
	}

	// Если ни одного заголовка не было найдено, отдаем пустой массив
	if len(titlesAndUrl) == 0 {
		return make([]models.ResponseTitlesAndUrls, 0), nil
	}

	return titlesAndUrl, nil
}

// Метод поиска заголовка в источнике
func getTitleFromSource(data *goquery.Document, url string) (models.ResponseTitlesAndUrls, error) {
	var responseData models.ResponseTitlesAndUrls
	findTitle := data.Find("title").Text()
	if findTitle == "" {
		return responseData, errors.New("title not found")
	}
	responseData.Title = findTitle
	responseData.Url = url

	return responseData, nil
}
