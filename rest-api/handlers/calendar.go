package handlers

import (
	"charg-go-microservice/rest-api/data"
	"log"
	"net/http"
)

type Calendar struct {
	l *log.Logger
}

func NewCalendar(l *log.Logger) *Calendar {
	return &Calendar{l}
}

func (c *Calendar) GetCalendar(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET Calendar")
	cal := data.GetCalendar()
	err := cal.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusInternalServerError)
	}
}
