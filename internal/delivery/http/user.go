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

	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/get_user_info/{userID}", handler.GetUserInfo).Methods("GET")
	router.HandleFunc("/update_user/{userID}", handler.UpdateUserInfo).Methods("PUT")
	router.HandleFunc("/get_avatar/{userID}", handler.GetUserAvatar).Methods("GET")
	router.HandleFunc("/get_pet_list/{userID}", handler.GetUsersPets).Methods("GET")
	router.HandleFunc("/add_pet/{userID}", handler.AddPet).Methods("POST")
	router.HandleFunc("/delete_pet", handler.DeletePet).Methods("DELETE")
	router.HandleFunc("/update_pet", handler.UpdatePet).Methods("PUT")
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, INVALID_BODY, http.StatusBadRequest)
		return
	}

	newUser := new(domain.ApiUserInfo)
	err = json.Unmarshal(body, newUser)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, INVALID_BODY, http.StatusBadRequest)
		return
	}

	loginResp, err := h.userUsecase.AddUser(newUser)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	jsonLoginResp, _ := json.Marshal(loginResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonLoginResp)
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

	loginResp, err := h.userUsecase.Login(loginInfo)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusForbidden)
		return
	}

	jsonLoginInfo, _ := json.Marshal(loginResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonLoginInfo)
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

func (h *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
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

	updInfo := new(domain.ApiUserUpdate)
	err = json.Unmarshal(body, updInfo)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, BAD_JSON_FORMAT, http.StatusBadRequest)
		return
	}

	err = h.userUsecase.UpdateUser(userID, updInfo)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	petIDToSend, err := h.userUsecase.AddPet(userID, newPet)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	jsonPetID, _ := json.Marshal(petIDToSend)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPetID)
}

func (h *UserHandler) DeletePet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := q.Get("userID")
	petID := q.Get("petID")

	if userID == "" || petID == "" {
		_ = responseTemplates.SendErrorMessage(w, BAD_QUERY_PARAMETERS, http.StatusBadRequest)
		return
	}

	err := h.userUsecase.DeletePet(userID, petID)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := q.Get("userID")
	petID := q.Get("petID")

	if userID == "" || petID == "" {
		_ = responseTemplates.SendErrorMessage(w, BAD_QUERY_PARAMETERS, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, INVALID_BODY, http.StatusBadRequest)
		return
	}

	updInfo := new(domain.ApiPetUpdate)
	err = json.Unmarshal(body, updInfo)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, BAD_JSON_FORMAT, http.StatusBadRequest)
		return
	}

	err = h.userUsecase.UpdatePet(userID, petID, updInfo)
	if err != nil {
		_ = responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
