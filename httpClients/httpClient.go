package httpClients

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
	"net/http"
)

type HttpClient struct {
	logger log.Logger
}

func NewHttpClient(logger log.Logger) HttpClientInterface {
	return &HttpClient{
		logger: logger,
	}
}

// ParseHtmlFromUrlSource Метод отправки запроса и возврата данных для парсинга
func (c *HttpClient) ParseHtmlFromUrlSource(url string) (*goquery.Document, error) {
	c.logger.Log("method", "ParseHtmlFromUrlSource")

	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
