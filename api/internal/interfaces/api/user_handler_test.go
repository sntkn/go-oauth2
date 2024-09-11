package api_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/config"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/api"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/response"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {

	cfg, err := config.GetEnv()
	if err != nil {
		log.Fatal("could not get env:", err)
	}

	dbConfig := &db.DBConfig{
		Host:     cfg.DBHost,
		Port:     uint16(cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	}

	database, err := db.Setup(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// テストケースの設定
	testID := uuid.New()
	//testUser := &user.User{
	//	ID:    testID,
	//	Name:  "テストユーザー",
	//	Email: "test@example.com",
	//}

	// エコーインスタンスの設定
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/"+testID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Bind params
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(testID.String())

	// インジェクションの設定
	i := &interfaces.Injections{
		DB: database, // DBは使用しないのでnilでOK
	}

	// ハンドラーの実行
	handler := api.GetUser(i)
	if assert.NoError(t, handler(c)) {

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response response.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Failed to retrieve users", response.Message)
		assert.Nil(t, response.Data)

		//assert.Equal(t, testID, response.ID)
		//assert.Equal(t, "テストユーザー", response.Name)
		//assert.Equal(t, "test@example.com", response.Email)
	}
}
