package main

import (
	"biscotti/handler"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	loggerStorage := handler.GetLoggerStorage()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func(storage *handler.LoggerStorage) {
		sig := <-sigs
		err := storage.Flush()
		fmt.Printf("exiting with %s signal\n", sig)
		if err != nil {
			fmt.Println(err)
		}
		done <- true
	}(loggerStorage)

	r := router.New()
	r.GET("/pixel", handler.PixelHandler)
	r.GET("/match", handler.MatchHandler)

	go func() {
		err := fasthttp.ListenAndServe(":3000", r.Handler)
		if err != nil {
			return
		}
	}()
	<-done
	fmt.Println("exited.")
}
