package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mgutz/logxi/v1"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
	"github.com/binarycode/trewoga/pkg/version"
)

type versionResponse struct {
	Version string `json:"version"`
}

type maintenanceResponse struct {
	Maintenance bool `json:"maintenance"`
}

const serviceContextKey = "service"

func Start() {
	router := mux.NewRouter()

	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/version", getVersion).Methods(http.MethodGet)
	v1.HandleFunc("/services/{token}/ping", middleware(ping)).Methods(http.MethodPost)
	v1.HandleFunc("/services/{token}/maintenance", middleware(enableMaintenance)).Methods(http.MethodPost)
	v1.HandleFunc("/services/{token}/maintenance", middleware(disableMaintenance)).Methods(http.MethodDelete)

	http.Handle("/", router)
}

func middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		if service, err := db.GetService(model.Service{Token: token}); err == nil {
			ctx := context.WithValue(r.Context(), serviceContextKey, service)
			handler.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.NotFound(w, r)
		}
	}
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	response := versionResponse{
		Version: version.Version,
	}
	write(w, response)
}

func ping(w http.ResponseWriter, r *http.Request) {
	s := getService(r)
	s.Ping()
	if err := db.SaveService(&s); err != nil {
		log.Fatal("Unable to save service", "err", err)
	}
}

func enableMaintenance(w http.ResponseWriter, r *http.Request) {
	s := getService(r)
	log.Info("Service maintenance", "service", s.Name)
	s.EnableMaintenance()
	if err := db.SaveService(&s); err != nil {
		log.Fatal("Unable to save service", "err", err)
	}
}

func disableMaintenance(w http.ResponseWriter, r *http.Request) {
	s := getService(r)
	log.Info("Service maintenance ended", "service", s.Name)
	s.DisableMaintenance()
	if err := db.SaveService(&s); err != nil {
		log.Fatal("Unable to save service", "err", err)
	}
}

func getService(r *http.Request) model.Service {
	return r.Context().Value(serviceContextKey).(model.Service)
}

func write(w http.ResponseWriter, response interface{}) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error encoding response", "err", err)
	}
}
