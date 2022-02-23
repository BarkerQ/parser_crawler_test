package httpClients

import "github.com/PuerkitoBio/goquery"

// HttpClientInterface Интерфейс для работы с http клиентами
type HttpClientInterface interface {
	ParseHtmlFromUrlSource(url string) (*goquery.Document, error)
}
