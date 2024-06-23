package main

import "testing"

func TestPing(t *testing.T) {
	_, err := ping()
	if err != nil {
		t.Fatal(err)
	}
}
