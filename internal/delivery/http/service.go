package http

import (
	"encoding/json"
	"io"
	"mainService/internal/domain"
	"mainService/internal/usecase"
	"mainService/pkg/responseTemplates"
	"net/http"

	"github.com/gorilla/mux"
)

type ServiceHandler struct {
	serviceUsecase usecase.IServiceUsecase
}

func NewServiceHandler(router *mux.Router, serviceUCase usecase.IServiceUsecase) {
	handler := &ServiceHandler{
		serviceUsecase: serviceUCase,
	}

	router.HandleFunc("/add_service/{userID}", handler.AddService).Methods("POST")
	router.HandleFunc("/get_service/{serviceID}", handler.GetService).Methods("GET")
	router.HandleFunc("/get_user_services/{userID}", handler.GetUserServices).Methods("GET")
	router.HandleFunc("/delete_service", handler.DeleteService).Methods("DELETE")
}

func (h *ServiceHandler) AddService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok {
		_ = responseTemplates.SendErrorMessage(w, BAD_GET_PARAMETER, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, INVALID_BODY, http.StatusBadRequest)
		return
	}

	newService := new(domain.ApiService)
	err = json.Unmarshal(body, newService)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, INVALID_BODY, http.StatusBadRequest)
		return
	}

	err = h.serviceUsecase.AddService(userID, newService)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ServiceHandler) GetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID, ok := vars["serviceID"]
	if !ok {
		_ = responseTemplates.SendErrorMessage(w, BAD_GET_PARAMETER, http.StatusBadRequest)
		return
	}

	service, err := h.serviceUsecase.GetServiceByID(serviceID)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	jsonServiceInfo, _ := json.Marshal(service)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonServiceInfo)
}

func (h *ServiceHandler) GetUserServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok {
		_ = responseTemplates.SendErrorMessage(w, BAD_GET_PARAMETER, http.StatusBadRequest)
		return
	}

	services, err := h.serviceUsecase.GetUserServices(userID)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	jsonServicesInfo, _ := json.Marshal(services)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonServicesInfo)
}

func (h *ServiceHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := q.Get("userID")
	serviceID := q.Get("serviceID")

	if userID == "" || serviceID == "" {
		_ = responseTemplates.SendErrorMessage(w, BAD_QUERY_PARAMETERS, http.StatusBadRequest)
		return
	}

	err := h.serviceUsecase.DeleteService(userID, serviceID)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
