package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthCodeInput struct {
	Code string `form:"code"`
}

type TokenRequest struct {
	Code      string `json:"code"`
	GrantType string `json:"grant_type"`
}

type AuthCodeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.GET("/callback", func(c *gin.Context) {
		var input AuthCodeInput
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		// POSTリクエストを送信するURL
		url := "http://localhost:8080/token"

		// POSTデータを作成
		reqData := TokenRequest{
			Code:      input.Code,
			GrantType: "authorization_code",
		}

		postData, err := json.Marshal(reqData)
		if err != nil {
			fmt.Println("Could not marshal json:", err)
			return
		}

		// POSTリクエストを作成
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
		if err != nil {
			fmt.Println("リクエストの作成エラー:", err)
			return
		}

		// リクエストヘッダーを設定（必要に応じて）
		req.Header.Set("Content-Type", "application/json")

		// HTTPクライアントを作成
		client := &http.Client{}

		// POSTリクエストを送信
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("リクエストの送信エラー:", err)
			return
		}
		defer resp.Body.Close()
		// レスポンスを読み取り
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("レスポンスの読み取りエラー:", err)
			return
		}
		var d AuthCodeResponse
		err = json.Unmarshal(body, &d)
		if err != nil {
			fmt.Println("Could not unmarshal auth code response:", err)
			return
		}
		c.JSON(200, d)
	})
	r.Run(":8000")
}
