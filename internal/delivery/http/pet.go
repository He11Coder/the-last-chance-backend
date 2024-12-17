package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"mainService/internal/usecase"
	"mainService/pkg/responseTemplates"
)

type PetHandler struct {
	petUsecase usecase.IPetUsecase
}

func NewPetHandler(router *mux.Router, petUCase usecase.IPetUsecase) {
	handler := &PetHandler{
		petUsecase: petUCase,
	}

	router.HandleFunc("/pet_info/{petID}", handler.GetPetInfo).Methods("GET")
	router.HandleFunc("/get_advice", handler.GetPetCareAdvice).Methods("GET")
}

func (h *PetHandler) GetPetInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	petID, ok := vars["petID"]
	if !ok {
		_ = responseTemplates.SendErrorMessage(w, BAD_GET_PARAMETER, http.StatusBadRequest)
		return
	}

	petInfo, err := h.petUsecase.GetPetInfo(petID)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	jsonInfo, _ := json.Marshal(petInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonInfo)
}

func (h *PetHandler) GetPetCareAdvice(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	animal := q.Get("animal")
	prompt := q.Get("prompt")

	if animal == "" || prompt == "" {
		_ = responseTemplates.SendErrorMessage(w, BAD_QUERY_PARAMETERS, http.StatusBadRequest)
		return
	}

	adviceURL := "http://127.0.0.1:8000/get_advice"

	params := url.Values{}
	params.Add("animal", animal)
	params.Add("prompt", prompt)

	fullAdviceURL := fmt.Sprintf("%s?%s", adviceURL, params.Encode())

	resp, err := http.Get(fullAdviceURL)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
