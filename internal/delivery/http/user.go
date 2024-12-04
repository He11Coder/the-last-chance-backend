package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"mainService/internal/domain"
	"mainService/internal/usecase"
	"mainService/pkg/responseTemplates"
)

type UserHandler struct {
	userUsecase usecase.IUserUsecase
}

func NewUserHandler(router *mux.Router, userUCase usecase.IUserUsecase) {
	handler := &UserHandler{
		userUsecase: userUCase,
	}

	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/get_user_info/{userID}", handler.GetUserInfo).Methods("GET")
	router.HandleFunc("/get_avatar/{userID}", handler.GetUserAvatar).Methods("GET")
	router.HandleFunc("/get_pet_list/{userID}", handler.GetUsersPets).Methods("GET")
	router.HandleFunc("/add_pet/{userID}", handler.AddPet).Methods("POST")
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, INVALID_BODY, http.StatusBadRequest)
		return
	}

	loginInfo := new(domain.LoginCredentials)
	err = json.Unmarshal(body, loginInfo)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, BAD_JSON_FORMAT, http.StatusBadRequest)
		return
	}

	//Check the credentials in database. Authorization.
	sessionID, err := h.userUsecase.Login(loginInfo)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, AUTH_ERROR, http.StatusForbidden)
		return
	}

	userInfo := domain.ApiUserInfo{UserID: sessionID}
	jsonUserInfo, _ := json.Marshal(userInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUserInfo)
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok {
		_ = responseTemplates.SendErrorMessage(w, MISSING_USER_ID, http.StatusBadRequest)
		return
	}

	userInfo, err := h.userUsecase.GetUserInfo(userID)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	jsonUserInfo, _ := json.Marshal(userInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUserInfo)
}

func (h *UserHandler) GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok {
		_ = responseTemplates.SendErrorMessage(w, BAD_GET_PARAMETER, http.StatusBadRequest)
		return
	}

	base64Image, err := h.userUsecase.GetUserAvatar(userID)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write([]byte(base64Image))
}

func (h *UserHandler) GetUsersPets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok {
		_ = responseTemplates.SendErrorMessage(w, BAD_GET_PARAMETER, http.StatusBadRequest)
		return
	}

	petIDs, err := h.userUsecase.GetUserPets(userID)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	jsonPetIDs, _ := json.Marshal(petIDs)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPetIDs)
}

func (h *UserHandler) AddPet(w http.ResponseWriter, r *http.Request) {
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

	newPet := new(domain.ApiPetInfo)
	err = json.Unmarshal(body, newPet)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, INVALID_BODY, http.StatusBadRequest)
		return
	}

	err = h.userUsecase.AddPet(userID, newPet)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
