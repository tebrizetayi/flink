package httpapi

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tebrizetayi/flink/app"
)

// API returns a handler for a set of routes.
func API(shutdown chan os.Signal, ttl int) http.Handler {

	app := app.NewApp(ttl)

	p := NewController(app)

	m := mux.NewRouter()

	m.HandleFunc("/location/{order_id}", p.ReadLocations).Methods("GET")
	m.HandleFunc("/location/{order_id}/now", p.SaveLocation).Methods("POST")
	m.HandleFunc("/location/{order_id}", p.DeleteLocation).Methods("DELETE")

	return m
}
