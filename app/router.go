package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"mainService/configs"
	deliveryHTTP "mainService/internal/delivery/http"
	"mainService/internal/repository/mongoTLC"
	"mainService/internal/repository/redisTLC"
	"mainService/internal/usecase"
	"mainService/pkg/swearWordsDetector"
)

func Run() error {
	if err := godotenv.Load("configs/.env"); err != nil {
		return err
	}

	configs.InitConfigs()

	err := swearWordsDetector.BuildAndCompileRegexp()
	if err != nil {
		return err
	}

	client, err := GetMongo()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	redisDB := GetRedis()
	defer redisDB.Close()

	db, err := InitDBAndIndexes(client)
	if err != nil {
		return err
	}

	userRepo := mongoTLC.NewMongoUserRepository(db)
	petRepo := mongoTLC.NewMongoPetRepository(db)
	serviceRepo := mongoTLC.NewMongoServiceRepository(db)
	sessionRepo := redisTLC.NewRedisAuthRepository(redisDB)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)
	petUsecase := usecase.NewPetUsecase(petRepo)
	serviceUsecase := usecase.NewServiceUsecase(serviceRepo, userRepo, petRepo)

	router := mux.NewRouter()
	deliveryHTTP.NewUserHandler(router, userUsecase)
	deliveryHTTP.NewPetHandler(router, petUsecase)
	deliveryHTTP.NewServiceHandler(router, serviceUsecase)

	http.Handle("/", router)

	fmt.Printf("\tstarting server at %s\n", ":8081")

	err = http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		return err
	}

	return nil
}
