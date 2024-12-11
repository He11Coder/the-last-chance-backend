package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"mainService/configs"
	deliveryHTTP "mainService/internal/delivery/http"
	"mainService/internal/repository/mongoTLC"
	"mainService/internal/repository/redisTLC"
	"mainService/internal/usecase"
)

/*type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	UserID    bson.ObjectID `json:"user_id,omitempty" bson:"_id,omitempty"`
	Username  string        `json:"username" bson:"name"`
	Contacts  string        `json:"contacts" bson:"contact"`
	UserImage string        `json:"user_image" bson:"avatar_url"`
}

type ErrorToSend struct {
	Message string `json:"message"`
}

var PetDB = map[int]domain.ApiPetInfo{
	1: {
		TypeOfAnimal: "Кошка",
		Name:         "Мурка",
		Info:         "Очень ласковая кошка. Имеет аллергию на мышей.",
		PetAvatar:    "cat1.png",
	},
	2: {
		TypeOfAnimal: "Кошка",
		Name:         "Маркиза",
		Info:         "Дико злая кошка. Имеет аллергию на людей.",
		PetAvatar:    "cat2.png",
	},
	3: {
		TypeOfAnimal: "Кот",
		Name:         "Гладиатор",
		Info:         "Бешеный кот. Ласков когда спит. Плохо пахнет",
		PetAvatar:    "cat3.png",
	},
	4: {
		TypeOfAnimal: "Змея",
		Name:         "Елизавета",
		Info:         "Почти не кусает. Кусает очень редко, но всегда с летальным исходом. Быть осторожным. Но вообще она очень милая и ласковая.",
		PetAvatar:    "snake.png",
	},
	5: {
		TypeOfAnimal: "Черепаха",
		Name:         "Танк",
		Info:         "Любит молоко...",
		PetAvatar:    "turtle.png",
	},
	6: {
		TypeOfAnimal: "Собака",
		Name:         "Кузя",
		Info:         "Собака с развитым любопытством.",
		PetAvatar:    "dog1.png",
	},
	7: {
		TypeOfAnimal: "Собака",
		Name:         "Лизка",
		Info:         "У собаки аллергия на других собак. Любит бегать.",
		PetAvatar:    "dog2.png",
	},
}

func GetPetInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	petID, convErr := strconv.Atoi(vars["petID"])
	if convErr != nil {
		errToSend := ErrorToSend{Message: "incorrect pet ID"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
		return
	}

	petInfo, ok := PetDB[petID]
	if !ok {
		errToSend := ErrorToSend{Message: "pet ID not found"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonErr)
		return
	}

	CURR_DIR, _ := os.Getwd()

	fileBytes, err := os.ReadFile(CURR_DIR + "/assets/pets" + "/" + petInfo.PetAvatar)
	if err != nil {
		errToSend := ErrorToSend{Message: "error while reading pet's avatar"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErr)
		return
	}

	base64Image := base64.StdEncoding.EncodeToString(fileBytes)
	petInfo.PetAvatar = base64Image

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
		return
	}

	var loginInfo LoginCredentials
	err = json.Unmarshal(body, &loginInfo)
	if err != nil {
		errToSend := ErrorToSend{Message: "invalid json format: must be with fields 'username' and 'password'"}
		jsonErr, _ := json.Marshal(errToSend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErr)
		return
	}

	//Check the credentials in database. Authorization.

	userInfo := UserInfo{UserID: bson.NewObjectID()}
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

func GetUserAvatar(w http.ResponseWriter, r *http.Request) {
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

func GetUsersPets(w http.ResponseWriter, r *http.Request) {
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
	petList := domain.PetIDList{PetIDs: petIDs}

	jsonPetList, _ := json.Marshal(petList)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPetList)
}*/

func Run() error {
	db, err := GetMongo()
	if err != nil {
		return err
	}
	defer db.Disconnect(context.TODO())

	redisDB := GetRedis()
	defer redisDB.Close()

	/*c := db.Database("tlc").Collection("user")

	u := UserInfo{
		Username:  "Сергей Иванов",
		Contacts:  "+79831238497",
		UserImage: "/assets/avatars/sergeant.png",
	}

	res, err := c.InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}
	fmt.Println("Ins ID:", res.InsertedID)*/

	//_ = mongoTLC.NewMongoUserRepository(db.Database("tlc"))

	//userObjID, err := primitive.ObjectIDFromHex("673de035bbe21f94ce848f4c")
	//userObjID, err := bson.ObjectIDFromHex("673de035bbe21f94ce848f4c")
	//if err != nil {
	//	return err
	//}
	//userRepo.AddPet(userObjID, &PetDB[2])

	userRepo := mongoTLC.NewMongoUserRepository(db.Database("tlc"))
	petRepo := mongoTLC.NewMongoPetRepository(db.Database("tlc"))
	sessionRepo := redisTLC.NewRedisAuthRepository(redisDB)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)
	petUsecase := usecase.NewPetUsecase(petRepo)

	router := mux.NewRouter()
	deliveryHTTP.NewUserHandler(router, userUsecase)
	deliveryHTTP.NewPetHandler(router, petUsecase)

	/*router.HandleFunc("/pet_info/{petID}", GetPetInfo).Methods("GET")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/get_user_info/{userID}", GetUserInfo).Methods("GET")
	router.HandleFunc("/get_avatar/{userID}", GetUserAvatar).Methods("GET")
	router.HandleFunc("/get_pet_list/{userID}", GetUsersPets).Methods("GET")*/

	http.Handle("/", router)

	fmt.Printf("\tstarting server at %s\n", ":8081")

	err = http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		return err
	}

	return nil
}
