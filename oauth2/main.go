package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Ginルーターの初期化
	r := gin.Default()

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
		// JSONレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// サーバーをポート8080で起動
	r.Run(":8080")
}
