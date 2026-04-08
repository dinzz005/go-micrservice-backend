package handler

import (
	"fmt"
	"log"
	"net/http"
)

type HealthHandler struct {
	 l *log.Logger
}

func NewHealthHandler(l *log.Logger) *HealthHandler{
	return &HealthHandler{l}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Health Check hit")
	fmt.Fprintln(w, "OK")
}

