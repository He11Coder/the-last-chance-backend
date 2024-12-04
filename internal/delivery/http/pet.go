package http

import (
	"encoding/json"
	"fmt"
	"net/http"

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
