package cfscan

import (
	"fmt"
	"testing"
)

func TestNextRunTime(t *testing.T) {
	time, err := NextRunTime("*/5 * * * *")
	if err != nil {
		panic(err)
	}
	fmt.Println(time)
}

func TestPreviousRunTime(t *testing.T) {
	time, err := PreviousRunTime("*/5 * * * *")
	if err != nil {
		panic(err)
	}
	fmt.Println(time)
}
