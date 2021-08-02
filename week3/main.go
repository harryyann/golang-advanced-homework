package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){
	serverChan := make(chan struct{})
	signalChan := make(chan struct{})
	signalCtx, cancel := context.WithCancel(context.Background())

	eg, _ := errgroup.WithContext(context.Background())

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello, world"))
		if err != nil{
			log.Println(err)
		}
	})

	// simulation server exit
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		serverChan <- struct{}{}
	})

	// start http server listening
	eg.Go(func() error {
		log.Println("server running...")
		err := server.ListenAndServe()
		if err != nil{
			serverChan <- struct{}{}
		}
		return err
	})

	// start signal listening
	eg.Go(func() error {
		log.Println("signal listening...")
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
		for{
			select {
			case <- c:
				signalChan <- struct{}{}
				return nil
			case <- signalCtx.Done():
				// stop signal listening
				return nil
			}
		}
	})

	select {
	case <- serverChan:
		// stop signal listening
		cancel()
		log.Fatal("server running error, stop signal listening and exit")
	case <- signalChan:
		log.Println("got interrupt signal, stop http server and exit")
		// shut down server
		ctx, _ := context.WithTimeout(context.Background(), time.Second * 10)
		err := server.Shutdown(ctx)
		if err != nil{
			log.Fatal("stop server err: ", err.Error())
		}
	}
}