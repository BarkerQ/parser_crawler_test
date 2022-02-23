package models

func NewResponse(code int, text string) Response {
	return Response{
		Code: code,
		Body: text,
	}
}

type Response struct {
	Code int    `json:"code"`
	Body string `json:"body"`
}

type UrlsRequest struct {
	Urls []string `json:"urls"`
}

type ResponseTitlesAndUrls struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}
