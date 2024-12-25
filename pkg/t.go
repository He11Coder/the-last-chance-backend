package main

import (
	"encoding/base64"
	"fmt"
	"mainService/pkg/nsfwFilter"
	"os"
)

func main() {
	fileBytes1, err := os.ReadFile("assets/avatars/nude_ass.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	fileBytes2, err := os.ReadFile("assets/avatars/chick.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	base64String1 := base64.StdEncoding.EncodeToString(fileBytes1)
	base64String2 := base64.StdEncoding.EncodeToString(fileBytes2)

	results := nsfwFilter.RunInParallel(base64String1, base64String2)
	for _, res := range results {
		if res.ProcessingErr != nil {
			fmt.Printf("Processing error: %+v\n", err)
		}

		fmt.Printf("%+v\n\n", res.Inf)
	}

	/*res, err := nsfwFilter.IsSafeForWork(base64String)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)*/
}
