package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/chris-argirop/charg-go-microsrvice/rest-api/db"
	"github.com/chris-argirop/charg-go-microsrvice/rest-api/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var totalRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of http requests",
	},
	[]string{"path"},
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

		totalRequest.With(prometheus.Labels{"path": "/api/"}).Inc()
	})
}

func init() {
	prometheus.Register(totalRequest)
}

func main() {
	l := log.New(os.Stdout, "charg-api", log.LstdFlags)

	dbaseDriver := os.Getenv("DB_DRIVER")
	dsourceName := os.Getenv("DS_NAME")

	db, err := db.NewDatabase(dbaseDriver, dsourceName)
	defer db.CloseConnection()
	if err != nil {
		l.Fatal(err)
	}
	fmt.Println("Successfull Connection to Database!")

	eh := handlers.NewExpense(l, db)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/get", eh.GetExpenses)
	getRouter.HandleFunc("/api/get/{id:[0-9]+}", eh.GetExpense)
	getRouter.HandleFunc("/api/delete/{id:[0-9]+}", eh.DeleteExpense)
	getRouter.HandleFunc("/api/clearall", eh.ClearTable)
	getRouter.HandleFunc("/api/get/{vendor:(?s).*}", eh.GetExpensesByVendor)

	sm.Use(prometheusMiddleware)

	sm.Path("/metrics").Handler(promhttp.Handler())

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/api/update/{id:[0-9]+}", eh.UpdateExpenses)
	putRouter.Use(eh.MiddlewarevalidateExpense)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/api/add", eh.AddExpenses)
	postRouter.Use(eh.MiddlewarevalidateExpense)

	s := &http.Server{
		Addr:         ":9100",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// start the server
	go func() {
		l.Println("Starting server on port 9100")

		err := s.ListenAndServe()
		if err != nil {
			l.Println("Error starting server", "error", err)
			os.Exit(1)
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
