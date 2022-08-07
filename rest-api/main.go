package main

import (
	"charg-go-microservice/rest-api/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "charg-api", log.LstdFlags)

	eh := handlers.NewExpense(l)
	ch := handlers.NewCalendar(l)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", eh.GetExpenses)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/update/{id:[0-9]+}", eh.UpdateExpenses)
	putRouter.Use(eh.MiddlewarevalidateExpense)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/add", eh.AddExpenses)
	postRouter.Use(eh.MiddlewarevalidateExpense)

	calRouter := sm.Methods(http.MethodGet).Subrouter()
	calRouter.HandleFunc("/calendar", ch.GetCalendar)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown\n", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)
}
