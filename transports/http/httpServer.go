package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"os"
	"parserProject/models"
	"parserProject/service"
	"parserProject/transports"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Server struct {
	serverHttp      *http.Server
	parserInterface service.ServiceInterface
	logger          log.Logger
}

func NewServerHttp(parserInterface service.ServiceInterface, logger log.Logger) *Server {
	return &Server{
		parserInterface: parserInterface,
		logger:          logger,
	}
}

func (s *Server) handler(endpoints transports.Endpoints) http.Handler {
	// Добавляем кастомные валидации
	s.addRulesValidate()

	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(s.logger)),
		kithttp.ServerErrorEncoder(s.encodeErrorResponse),
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			return context.WithValue(context.Background(), "x-api-key", request.Header.Get("x-api-key"))
		}),
	}

	r.HandleFunc("/", s.getDefaultData())

	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("../docs/"))))

	r.Methods("POST").Path("/parse").Handler(kithttp.NewServer(
		endpoints.GetTitleFromUrls,
		s.decodeRequestUrls,
		s.encodeResponse,
		options...,
	))

	r.MethodNotAllowedHandler = s.errorHandlerNotAllowed()
	r.NotFoundHandler = s.errorHandlerNotFound()

	return r
}

func (s *Server) Run() error {
	logrus.Info("Start service")

	httpServer := &http.Server{
		Addr:    ":" + os.Getenv("SERVICE_PORT"),
		Handler: s.handler(transports.MakeEndpoints(s.parserInterface, os.Getenv("API_KEY"))),
	}
	s.serverHttp = httpServer

	err := httpServer.ListenAndServe()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Server) Shutdown() error {
	logrus.Info("Shutdown server...")

	// Задаем свой таймаут
	getTimeout := os.Getenv("TIMEOUT_HTTP_SERVER")
	atoi, _ := strconv.Atoi(getTimeout)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(atoi)*time.Second)
	err := s.serverHttp.Shutdown(ctx)
	return err
}

func (s *Server) encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok {
		s.encodeErrorResponse(ctx, e, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func (s *Server) getDefaultData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(models.NewResponse(200, "Please, choose method."))
	}
}

func (s *Server) encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := 500
	if err == transports.NoaAuth {
		code = 401
	} else if strings.Contains(err.Error(), "validationError") {
		code = 400
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.NewResponse(code, err.Error()))
}

func (s *Server) errorHandlerNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(models.NewResponse(405, transports.NotAllowed.Error()))
	})
}

func (s *Server) errorHandlerNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(models.NewResponse(404, transports.NotFound.Error()))
	})
}

func (s *Server) decodeRequestUrls(_ context.Context, r *http.Request) (request interface{}, err error) {
	var itemRequest models.UrlsRequest

	// Валидируем тело запроса
	rules := govalidator.MapData{
		"urls": []string{"required", "array"},
	}
	opts := govalidator.Options{
		Request: r,
		Data:    &itemRequest,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	if len(e) > 0 {
		return nil, formError(e)
	}

	return itemRequest, nil
}

func formError(message interface{}) error {
	errs := map[string]interface{}{"validationError": message}
	requestJson, err := json.Marshal(errs)
	if err != nil {
		return err
	}
	return errors.New(string(requestJson))
}

// Метод добавления кастомных ошибок валидации
func (s *Server) addRulesValidate() {
	// кастомная проверка на маасив строк
	govalidator.AddCustomRule("array", func(field string, rule string, message string, value interface{}) error {
		values := value.([]string)
		for _, c := range values {
			if !isString(c) {
				return errors.New("the field must be an array of string")
			}
		}
		return nil
	})
}

func isString(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return true
		}
	}
	return false
}
