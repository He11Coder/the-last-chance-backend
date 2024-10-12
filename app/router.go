package app

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type ErrorToSend struct {
	Message string `json:"message"`
}

func get_pet_info(w http.ResponseWriter, r *http.Request) {
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

func Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/pet_info/{petID}", get_pet_info).Methods("GET")

	http.Handle("/", router)

	fmt.Printf("\tstarting server at %s\n", ":8081")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		return err
	}

	return nil
}
