package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/config"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/model"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/query"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/api"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserExists(t *testing.T) {
	t.Parallel()
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
	tx := database.Begin()
	defer tx.Rollback()

	// テストケースの設定
	testID := uuid.New()
	testUser := &model.User{
		ID:        testID.String(),
		Name:      "テストユーザー",
		Email:     "test@example.com",
		Password:  "test1234",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	q := query.Use(tx)
	userQuery := q.User

	err = userQuery.WithContext(context.Background()).Create(testUser)
	require.NoError(t, err)

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
		DB: tx,
	}

	// ハンドラーの実行
	handler := api.NewHandler(i)
	if assert.NoError(t, handler.GetUser(c)) {

		assert.Equal(t, http.StatusOK, rec.Code)

		var response response.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "success", response.Status)

		data, ok := response.Data.(map[string]interface{})
		require.True(t, ok)

		fmt.Println(data)

		assert.Equal(t, testID.String(), data["id"])
		assert.Equal(t, "テストユーザー", data["name"])
		assert.Equal(t, "test@example.com", data["email"])
	}
}
func TestGetUserNotExists(t *testing.T) {
	t.Parallel()
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
	tx := database.Begin()
	defer tx.Rollback()

	// テストケースの設定
	testID := uuid.New()

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
		DB: tx,
	}

	// ハンドラーの実行
	handler := api.NewHandler(i)
	if assert.NoError(t, handler.GetUser(c)) {

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response response.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Failed to retrieve users", response.Message)
		assert.Nil(t, response.Data)
	}
}
