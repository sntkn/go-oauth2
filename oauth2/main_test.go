package main

import (
	"fmt"
	"testing"
)

func TestPing(t *testing.T) {
	_, err := ping()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("success")
}
