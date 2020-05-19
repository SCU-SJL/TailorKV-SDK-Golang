package main

import (
	"fmt"
	"log"
	"tailorSDK"
)

func main() {
	tailor, err := tailorSDK.Connect("localhost", "8448", "", "")
	if err != nil {
		log.Fatal(err)
	}
	defer tailor.Shutdown()

	err = tailor.Set("name", "sjl")
	if err != nil {
		log.Fatal(err)
	}

	val, err := tailor.Get("name", 32)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name = %s", val)
}
