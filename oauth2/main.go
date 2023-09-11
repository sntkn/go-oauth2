package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthorizeInput struct {
	ResponseType string `json:"response_type"`
	ClientId     string `json:"client_id"`
	Scope        string `json:"Scope"`
	RedirectURI  string `json:"redirect_uri"`
	State        string `json:"State"`
}

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// GETリクエストを受け取るエンドポイントの定義
	r.GET("/authorize", func(c *gin.Context) {
		// /authorize?response_type=code&client_id=abcdefg&Scope=read&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&State=ok
		input := AuthorizeInput{}

		input.ResponseType = c.Query("response_type")
		log.Printf("Response type: %s", c.Query("response_type"))
		if input.ResponseType == "" {
			c.HTML(http.StatusBadRequest, "Invalid response type", nil)
			return
		}

		input.ClientId = c.Query("client_id")
		if input.ClientId == "" {
			c.HTML(http.StatusBadRequest, "Invalid client_id", nil)
			return
		}

		input.Scope = c.Query("scope")
		if input.Scope == "" {
			c.HTML(http.StatusBadRequest, "Invalid scope", nil)
			return
		}

		input.RedirectURI = c.Query("redirect_uri")
		if input.RedirectURI == "" {
			c.HTML(http.StatusBadRequest, "Invalid redirect_uri", nil)
			return
		}

		input.State = c.Query("state")
		if input.State == "" {
			c.HTML(http.StatusBadRequest, "Invalid state", nil)
			return
		}
		log.Printf("%+v\n", input)
		c.HTML(http.StatusOK, "index.html", gin.H{"input": input})
	})

	// サーバーをポート8080で起動
	r.Run(":8080")
}
