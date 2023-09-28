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

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/callback", func(c *gin.Context) {
		var input AuthCodeInput
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		// POSTリクエストを送信するURL
		uri := "http://localhost:8080/token"

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
		req, err := http.NewRequest("POST", uri, bytes.NewBuffer(postData))
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

		c.SetCookie("access_token", d.AccessToken, 3600, "/", "localhost", false, true)
		c.SetCookie("refresh_token", d.RefreshToken, 3600, "/", "localhost", false, true)
		c.SetCookie("expiry", fmt.Sprintf("%d", d.Expiry), 3600, "/", "localhost", false, true)

		c.Redirect(http.StatusFound, "/home")
	})
	r.GET("/home", func(c *gin.Context) {
		// TODO: check cookie access token
		token, err := c.Cookie("access_token")
		if err != nil {
			fmt.Println("Error getting access token")
			c.Redirect(http.StatusFound, "/")
			return
		}

		// TODO: request user by token
		// GETリクエストを送信するURL
		uri := "http://localhost:8080/me"

		// GETリクエストを作成
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			fmt.Println("リクエストの作成エラー:", err)
			return
		}

		// リクエストヘッダーを設定（必要に応じて）
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// HTTPクライアントを作成
		client := &http.Client{}

		// GETリクエストを送信
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

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, string(body))
			return
		}

		var d UserResponse
		err = json.Unmarshal(body, &d)
		if err != nil {
			fmt.Println("Could not unmarshal auth code response:", err)
			return
		}
		c.HTML(http.StatusOK, "home.html", gin.H{"data": d})
	})

	r.GET("/logout", func(c *gin.Context) {
		// TODO: check cookie access token
		token, err := c.Cookie("access_token")
		if err != nil {
			fmt.Println("Error getting access token")
			c.Redirect(http.StatusFound, "/")
			return
		}

		// TODO: request user by token
		// DELETE リクエストを送信するURL
		uri := "http://localhost:8080/token"

		// DELETE リクエストを作成
		req, err := http.NewRequest("DELETE", uri, nil)
		if err != nil {
			fmt.Println("リクエストの作成エラー:", err)
			return
		}

		// リクエストヘッダーを設定（必要に応じて）
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// HTTPクライアントを作成
		client := &http.Client{}

		// GETリクエストを送信
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("リクエストの送信エラー:", err)
			c.HTML(http.StatusUnauthorized, "400.html", nil)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.HTML(http.StatusUnauthorized, "400.html", nil)
			return
		}

		c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
		c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
		c.SetCookie("expiry", "", -1, "/", "localhost", false, true)

		c.HTML(http.StatusOK, "logout.html", nil)
	})

	r.Run(":8000")
}
