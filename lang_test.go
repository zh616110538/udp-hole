package main

import (
	"fmt"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	var m = make(map[string]bool)
	if !m["123"] {
		fmt.Println(m["1"])
	}
}

func TestChannal(t *testing.T) {
	stop := make(chan bool)
	go func() {
		run := true
		for run {
			select {
			case run = <-stop:
			default:
				fmt.Println("no data")
			}
		}
	}()
	time.Sleep(time.Millisecond * 10)
	stop <- true
}
