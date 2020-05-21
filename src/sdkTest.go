package main

import (
	"fmt"
	"log"
	"tailorSDK"
	"time"
)

func main() {
	tailor, err := tailorSDK.Connect("127.0.0.1", "8448", "123456", "SJL *loves* ZHH-")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tailor.Ttl("me")
	fmt.Println(err)
	tailor.Setex("me", "sjl", 10*time.Second)
	time.Sleep(1 * time.Second)
	fmt.Println(tailor.Ttl("me"))
	time.Sleep(1 * time.Second)
	fmt.Println(tailor.Ttl("me"))
	time.Sleep(1 * time.Second)
	fmt.Println(tailor.Ttl("me"))
	time.Sleep(1 * time.Second)
	fmt.Println(tailor.Ttl("me"))
	time.Sleep(1 * time.Second)
	fmt.Println(tailor.Ttl("me"))
}
