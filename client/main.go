package main

import "github.com/gin-gonic/gin"

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.GET("/callback", func(c *gin.Context) {
		c.JSON(200, gin.H{"result": "OK"})
	})
	r.Run(":8000")
}
