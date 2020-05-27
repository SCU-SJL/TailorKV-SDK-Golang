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

func TestTailor_Ttl(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	_, err := tailor.Ttl("123")
	if err != nil {
		fmt.Println(err)
	}
	_ = tailor.Setex("10s", "10 seconds", 10*time.Second)
	ttl, err := tailor.Ttl("10s")
	if err != nil {
		t.Errorf("Ttl() err: %v", err)
	} else {
		fmt.Println("Ttl() res: ", ttl)
	}

	ttl, err = tailor.Ttl("10s")
	if err != nil {
		t.Errorf("Ttl() err: %v", err)
	} else {
		fmt.Println("Ttl() res: ", ttl)
	}
}

func TestTailor_Unlink(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	err := tailor.Set("me", "sjl")
	if err != nil {
		t.Error(err)
	}
	err = tailor.Unlink("me")
	if err != nil {
		t.Error(err)
	}
}

func TestTailor_Incr(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	err := tailor.Set("num1", "1")
	if err != nil {
		t.Error(err)
	}
	err = tailor.Incr("num1")
	if err != nil {
		t.Error(err)
	}
	err = tailor.Incrby("num1", 10)
	if err != nil {
		t.Error(err)
	}
}

func TestTailor_Keys(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	ks, err := tailor.Keys("[a-z]+", 1024)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ks)
}

func TestTailor_Save(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	err := tailor.Set("me", "sjl")
	if err != nil {
		t.Fatal(err)
	}
	err = tailor.Save("testSDK.tkv")
	if err != nil {
		t.Error(err)
	}
}

func TestTailor_Load(t *testing.T) {
	if tErr != nil {
		t.Fatal(tErr)
	}
	err := tailor.Load("testSDK.tkv")
	if err != nil {
		t.Error(err)
	}
}
