package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Ginルーターの初期化
	r := gin.Default()

	// GETリクエストを受け取るエンドポイントの定義
	r.GET("/api/hello", func(c *gin.Context) {
		// JSONレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// サーバーをポート8080で起動
	r.Run(":8080")
}
