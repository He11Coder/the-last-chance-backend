package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type PetInfo struct {
	TypeOfAnimal string   `json:"type_of_animal"`
	Breed        string   `json:"breed"`
	Name         string   `json:"name"`
	Allergens    []string `json:"allergens"`
	Preferences  []string `json:"preferences"`
}

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Contacts  string `json:"contacts"`
	UserImage string `json:"user_image"`
}

type ErrorToSend struct {
	Message string `json:"message"`
}

func GetPetInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, convErr := strconv.Atoi(vars["petID"])
	if convErr != nil {
		errToSend := ErrorToSend{Message: "incorrect pet ID"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
	}

	petInfo := PetInfo{
		TypeOfAnimal: "Собака",
		Breed:        "Корги",
		Name:         "Мага",
		Allergens:    []string{"Конфеты", "Молоко"},
		Preferences:  []string{"Кости", "Мясо"},
	}

	jsonInfo, _ := json.Marshal(petInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonInfo)
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errToSend := ErrorToSend{Message: "invalid request body"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
	}

	var loginInfo LoginCredentials
	err = json.Unmarshal(body, &loginInfo)
	if err != nil {
		errToSend := ErrorToSend{Message: "invalid json format: must be with fields 'username' and 'password'"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
	}

	//Check the credentials in database. Authorization.

	userInfo := UserInfo{UserID: 1}
	jsonUserInfo, _ := json.Marshal(userInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUserInfo)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, convErr := strconv.Atoi(vars["userID"])
	if convErr != nil {
		errToSend := ErrorToSend{Message: "incorrect user ID"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
	}

	CURR_DIR, _ := os.Getwd()

	fileBytes, err := os.ReadFile(CURR_DIR + "/assets/avatars/sergeant.png")
	if err != nil {
		errToSend := ErrorToSend{Message: "error while reading user's avatar"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErr)
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

func GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, convErr := strconv.Atoi(vars["userID"])
	if convErr != nil {
		errToSend := ErrorToSend{Message: "incorrect user ID"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
	}

	CURR_DIR, _ := os.Getwd()

	fileBytes, err := os.ReadFile(CURR_DIR + "/assets/avatars/sergeant.png")
	if err != nil {
		errToSend := ErrorToSend{Message: "error while reading user's avatar"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErr)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/pet_info/{petID}", GetPetInfo).Methods("GET")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/get_user_info/{userID}", GetUserInfo).Methods("GET")
	router.HandleFunc("/get_avatar/{userID}", GetUserAvatar).Methods("GET")

	http.Handle("/", router)

	fmt.Printf("\tstarting server at %s\n", ":8081")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		return err
	}

	return nil
}
