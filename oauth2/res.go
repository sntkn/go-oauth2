package main

import (
	"encoding/json"
	"fmt"
)

type Res struct {
	Requested int `json:"requested"`
}

func (r *Res) SetRequested(i int) {
	r.Requested = i
}

type PostRes struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

type WithRequested[T any] struct {
	Requested int `json:"requested"`
	Data      T   `json:"data"`
}

func ResHandler[T any](data T) string {
	wrapped := WithRequested[T]{
		Requested: 1,
		Data:      data,
	}

	jsonData, _ := json.Marshal(wrapped)
	return string(jsonData)
}

func main() {
	res := PostRes{
		Posts: []Post{
			{ID: 1, Address: "Tokyo"},
			{ID: 2, Address: "Osaka"},
		},
	}
	fmt.Println(ResHandler(res))
}
