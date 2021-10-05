package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/tebrizetayi/flink/app"
)

var MIMEApplicationJSON = "application/json"

type ReadLocationResponse struct {
	OrderId string         `json:"order_id"`
	History []app.Location `json:"history"`
}

//Controller represents the API method handlers set.
type Controller struct {
	app app.App
}

//NewController creates new controller.
func NewController(app app.App) Controller {
	return Controller{app}
}

func (h Controller) SaveLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["order_id"]

	location := app.Location{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&location)

	if handleIfError(w, err, http.StatusBadRequest) {
		return
	}

	err = h.app.SaveLocation(r.Context(), location, orderId)

	if handleIfError(w, err, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h Controller) ReadLocations(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["order_id"]

	maxParam := ""
	if r.URL.Query()["max"] != nil && len(r.URL.Query()["max"]) > 0 {
		maxParam = r.URL.Query()["max"][0]
	}
	max, _ := strconv.Atoi(maxParam)

	locations, err := h.app.ReadLocations(r.Context(), orderId, max)

	if handleIfError(w, err, http.StatusInternalServerError) {
		return
	}
	resp := ReadLocationResponse{OrderId: orderId, History: locations}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", MIMEApplicationJSON)
	json.NewEncoder(w).Encode(resp)
}

func (h Controller) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderId := vars["order_id"]

	err := h.app.DeleteLocation(r.Context(), orderId)

	if handleIfError(w, err, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

type APIError struct {
	Error struct {
		// Message is a developer friendly representation of the error issue
		Message string `json:"message"`
	} `json:"error"`
}

func handleIfError(w http.ResponseWriter, err error, code int) bool {
	if err == nil {
		return false
	}
	apiErr := APIError{}
	apiErr.Error.Message = err.Error()

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		log.Println(`ERROR`, err.Error())
	}
	return true
}
