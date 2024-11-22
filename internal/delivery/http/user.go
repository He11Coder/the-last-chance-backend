package http

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

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
	_, ok := vars["userID"]
	if !ok {
		_ = responseTemplates.SendErrorMessage(w, MISSING_USER_ID, http.StatusBadRequest)
		return
	}

	CURR_DIR, _ := os.Getwd()

	fileBytes, err := os.ReadFile(CURR_DIR + "/assets/avatars/sergeant.png")
	if err != nil {
		errToSend := ErrorToSend{Message: "error while reading user's avatar"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErr)
		return
	}

	base64Image := base64.StdEncoding.EncodeToString(fileBytes)

	userInfo := UserInfo{
		Username:  "Сергей Иванов",
		Contacts:  "+79831238497",
		UserImage: base64Image,
	}

	jsonUserInfo, _ := json.Marshal(userInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUserInfo)
}

func (h *UserHandler) GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, convErr := strconv.Atoi(vars["userID"])
	if convErr != nil {
		errToSend := ErrorToSend{Message: "incorrect user ID"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
		return
	}

	CURR_DIR, _ := os.Getwd()

	fileBytes, err := os.ReadFile(CURR_DIR + "/assets/avatars/sergeant.png")
	if err != nil {
		errToSend := ErrorToSend{Message: "error while reading user's avatar"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func (h *UserHandler) GetUsersPets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, convErr := strconv.Atoi(vars["userID"])
	if convErr != nil {
		errToSend := ErrorToSend{Message: "incorrect user ID"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
		return
	}

	petIDs := []int{1, 2, 3, 4, 5, 6, 7}
	petList := PetIDList{PetIDs: petIDs}

	jsonPetList, _ := json.Marshal(petList)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPetList)
}
