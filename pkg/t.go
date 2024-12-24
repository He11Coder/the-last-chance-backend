package main

import (
	"encoding/base64"
	"fmt"
	"mainService/pkg/nsfwFilter"
	"os"
)

func main() {
	fileBytes, err := os.ReadFile("assets/avatars/nude_ass.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	base64String := base64.StdEncoding.EncodeToString(fileBytes)

	res, err := nsfwFilter.IsSafeForWork(base64String)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
}
