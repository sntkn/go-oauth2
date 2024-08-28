package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthorizationUsecase struct {
	mock.Mock
}

func (m *MockAuthorizationUsecase) Invoke(c *gin.Context, email string, password string) (string, error) {
	args := m.Called(c)
	return args.String(0), args.Error(1)
}

func TestAuthorizationHandler(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()
	r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

	r.Use(func(c *gin.Context) {
		c.Set("db", &repository.OAuth2RepositoryMock{
			FindUserByEmailFunc: func(email string) (repository.User, error) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
				return repository.User{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Name:      "client Name",
					Email:     "test@example.com",
					Password:  string(hashedPassword),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			RegisterOAuth2CodeFunc: func(c *repository.Code) error {
				return nil
			},
		})
		c.Set("session", &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			DelSessionDataFunc: func(c *gin.Context, key string) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*usecases.AuthorizeInput) = usecases.AuthorizeInput{
					ClientID:    "00000000-0000-0000-0000-000000000000",
					Scope:       "read",
					RedirectURI: "http://example.com",
				}
				return nil
			},
		})
		c.Set("cfg", &config.Config{})
		c.Next() // 次のミドルウェア/ハンドラへ
	})

	// authorizationハンドラをセット
	r.POST("/authorization", func(c *gin.Context) {
		AuthorizationHandler(c)
	})

	// 構造体からフォームデータを生成
	values := url.Values{}
	values.Set("email", "test@example.com")
	values.Add("password", "test1234")

	req, err := http.NewRequest(http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	assert.NoError(t, err)

	// レスポンスを記録するためのレスポンスレコーダを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	r.ServeHTTP(w, req)

	// ステータスコードが302 StatusFoundであることを確認
	assert.Equal(t, http.StatusFound, w.Code)
}

func TestAuthorization(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful authorization", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		// 構造体からフォームデータを生成
		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")
		c.Request = httptest.NewRequest(http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockAuthorizationUsecase)
		mockUC.On("Invoke", mock.Anything).Return("http://example.com/?code=abcdefg", nil)

		authorization(c, mockUC, s)
		c.Writer.WriteHeaderNow() // POST だとヘッダの書き込みが行われず、200を返してしまうのでここで書き込み https://stackoverflow.com/questions/76319196/unit-testing-of-gins-context-redirect-works-for-get-response-code-but-fails-for

		assert.Equal(t, "http://example.com/?code=abcdefg", w.Header().Get("Location"))
		assert.Equal(t, http.StatusFound, w.Code)

		mockUC.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		// 構造体からフォームデータを生成
		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "")
		c.Request = httptest.NewRequest(http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockAuthorizationUsecase)
		mockUC.On("Invoke", mock.Anything).Return("http://example.com/?code=abcdefg", nil)

		authorization(c, mockUC, s)
		c.Writer.WriteHeaderNow()

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/signin", w.Header().Get("Location"))
	})

	t.Run("usecase error(bad request)", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		// 構造体からフォームデータを生成
		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")
		c.Request = httptest.NewRequest(http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockAuthorizationUsecase)
		mockUC.On("Invoke", mock.Anything).Return("", &cerrs.UsecaseError{Code: http.StatusBadRequest})

		authorization(c, mockUC, s)
		c.Writer.WriteHeaderNow() // POST だとヘッダの書き込みが行われず、200を返してしまうのでここで書き込み https://stackoverflow.com/questions/76319196/unit-testing-of-gins-context-redirect-works-for-get-response-code-but-fails-for

		assert.Equal(t, "/signin", w.Header().Get("Location"))
		assert.Equal(t, http.StatusFound, w.Code)

		mockUC.AssertExpectations(t)
	})

	t.Run("usecase error(internal server error)", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		// 構造体からフォームデータを生成
		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")
		c.Request = httptest.NewRequest(http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockAuthorizationUsecase)
		mockUC.On("Invoke", mock.Anything).Return("", errors.New("internal error"))

		authorization(c, mockUC, s)
		c.Writer.WriteHeaderNow() // POST だとヘッダの書き込みが行われず、200を返してしまうのでここで書き込み https://stackoverflow.com/questions/76319196/unit-testing-of-gins-context-redirect-works-for-get-response-code-but-fails-for

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockUC.AssertExpectations(t)
	})
}
