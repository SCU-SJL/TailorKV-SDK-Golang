package tailorSDK

import (
	"fmt"
	"testing"
	"time"
)

var tailor, tErr = Connect("localhost", "8448", "123456", "SJL *loves* ZHH-")

func TestConnect(t *testing.T) {
	_, err := Connect("localhost", "8448", "123456", "SJL *loves* ZHH-")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTailor_Set(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	err := tailor.Set("me", "Jack")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTailor_Cnt(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	count, err := tailor.Cnt()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Cnt() =", count)
}

func TestTailor_Setex(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	err := tailor.Setex("she", "Zhh", 20*time.Second)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTailor_Setnx(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	err := tailor.Setnx("me", "SJL")
	if err == nil {
		t.Errorf("setnx failed")
	}
	fmt.Println(err)
}

func TestTailor_Del(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	err := tailor.Del("me")
	if err != nil {
		t.Errorf("func Del(): %v", err)
	}
}
