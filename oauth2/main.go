package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// GETリクエストを受け取るエンドポイントの定義
	r.GET("/authorize", func(c *gin.Context) {
		// Get parameters
		// response_type
		// client_id
		// scope
		// redirect_uri
		// state
		// validate parameters
		// redirect to authorization
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// サーバーをポート8080で起動
	r.Run(":8080")
}
