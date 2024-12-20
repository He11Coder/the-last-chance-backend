package main

import (
	"fmt"

	"mainService/app"
)

func main() {
	err := app.Run()

	if err != nil {
		fmt.Printf("err: %v\n in detail: %v\n", fmt.Errorf("the server encountered a problem while starting"), err)
		return
	}
}
