package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	http2 "net/http"
	"os"
	"os/signal"
	"parserProject/httpClients"
	"parserProject/logs"
	"parserProject/service"
	"parserProject/transports/http"
	"sync"
	"syscall"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	container.Provide(logs.NewLogger)
	container.Provide(httpClients.NewHttpClient)
	container.Provide(service.NewService)
	container.Provide(http.NewServerHttp)
	return container
}

func main() {
	container := BuildContainer()
	err := container.Invoke(func(server *http.Server) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			sigc := make(chan os.Signal)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			<-sigc

			err := server.Shutdown()
			if err != nil {
				if err == context.DeadlineExceeded {
					logrus.Println("shutdown: halted active connections")
				} else {
					logrus.Println(err)
				}
			}
			wg.Done()
		}()

		err := server.Run()
		if err != http2.ErrServerClosed {
			logrus.Fatalln(err)
		}

		wg.Wait()
	})

	if err != nil {
		logrus.Fatalln(err)
	}
}

// Метод инициализации переменного окружения
func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("No .env file found")
	}
}
