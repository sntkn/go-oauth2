package api_test

import (
	"context"
	"encoding/json"
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
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	cfg, err := config.GetEnv()
	require.NoError(t, err)

	dbConfig := &db.DBConfig{
		Host:     cfg.DBHost,
		Port:     uint16(cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	}

	database, err := db.Setup(dbConfig)
	require.NoError(t, err)

	tx := database.Begin()
	return tx, func() { tx.Rollback() }
}

func setupEchoContext(method, path string, params map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for key, value := range params {
		c.SetParamNames(key)
		c.SetParamValues(value)
	}

	return c, rec
}

func TestGetUserExists(t *testing.T) {
	t.Parallel()
	tx, cleanup := setupTestDB(t)
	defer cleanup()

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

	err := userQuery.WithContext(context.Background()).Create(testUser)
	require.NoError(t, err)

	c, rec := setupEchoContext(http.MethodGet, "/users/"+testID.String(), map[string]string{"id": testID.String()})

	i := &interfaces.Injections{DB: tx}
	handler := api.NewUserHandler(i)
	require.NoError(t, handler.GetUser(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var response response.Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "success", response.Status)

	data, ok := response.Data.(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, testID.String(), data["id"])
	assert.Equal(t, "テストユーザー", data["name"])
	assert.Equal(t, "test@example.com", data["email"])
}

func TestGetUserNotExists(t *testing.T) {
	t.Parallel()
	tx, cleanup := setupTestDB(t)
	defer cleanup()

	testID := uuid.New()

	c, rec := setupEchoContext(http.MethodGet, "/users/"+testID.String(), map[string]string{"id": testID.String()})

	i := &interfaces.Injections{DB: tx}
	handler := api.NewUserHandler(i)
	require.NoError(t, handler.GetUser(c))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response response.Response
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Failed to retrieve users", response.Message)
	assert.Nil(t, response.Data)
}
