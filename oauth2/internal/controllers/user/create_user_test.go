package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCreateUserUsecase struct {
	mock.Mock
}

func (m *MockCreateUserUsecase) Invoke(c *gin.Context, user repository.User) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestCreateUserHandler(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()
	r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

	r.Use(func(c *gin.Context) {
		c.Set("db", &repository.OAuth2RepositoryMock{
			ExistsUserByEmailFunc: func(email string) (bool, error) {
				return false, nil
			},
			CreateUserFunc: func(u *repository.User) error {
				return nil
			},
		})
		c.Set("session", &session.SessionClientMock{
			DelSessionDataFunc: func(c *gin.Context, key string) error {
				return nil
			},
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		})
		c.Set("cfg", &config.Config{})
		c.Next() // 次のミドルウェア/ハンドラへ
	})

	// サインインハンドラをセット
	r.POST("/signup", func(c *gin.Context) {
		CreateUserHandler(c)
	})

	// 構造体からフォームデータを生成
	values := url.Values{}
	values.Add("name", "test")
	values.Add("email", "test@example.com")
	values.Add("password", "test1234")

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/signup", strings.NewReader(values.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// レスポンスを記録するためのレスポンスレコーダを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	r.ServeHTTP(w, req)

	// ステータスコードが200 OKであることを確認
	assert.Equal(t, http.StatusFound, w.Code)
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful createUser", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		// 構造体からフォームデータを生成
		values := url.Values{}
		values.Add("name", "test")
		values.Add("email", "test@example.com")
		values.Add("password", "test1234")
		c.Request = httptest.NewRequest(http.MethodPost, "/singup", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockCreateUserUsecase)
		mockUC.On("Invoke", mock.Anything).Return(nil)

		createUser(c, mockUC, s)
		c.Writer.WriteHeaderNow() // POST だとヘッダの書き込みが行われず、200を返してしまうのでここで書き込み https://stackoverflow.com/questions/76319196/unit-testing-of-gins-context-redirect-works-for-get-response-code-but-fails-for

		assert.Equal(t, "/signup-finished", w.Header().Get("Location"))
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
		values.Add("name", "")
		values.Add("email", "test@example.com")
		values.Add("password", "test1234")
		c.Request = httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockCreateUserUsecase)
		mockUC.On("Invoke", mock.Anything).Return(nil)

		createUser(c, mockUC, s)
		c.Writer.WriteHeaderNow()

		assert.Equal(t, "/signup", w.Header().Get("Location"))
		assert.Equal(t, http.StatusFound, w.Code) // c.Bindだと400だが、c.ShouldBindだと302
	})

	t.Run("usecase error(bad request)", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		// 構造体からフォームデータを生成
		values := url.Values{}
		values.Add("name", "test")
		values.Add("email", "test@example.com")
		values.Add("password", "test1234")
		c.Request = httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockCreateUserUsecase)
		mockUC.On("Invoke", mock.Anything).Return(&cerrs.UsecaseError{Code: http.StatusBadRequest})

		createUser(c, mockUC, s)
		c.Writer.WriteHeaderNow() // POST だとヘッダの書き込みが行われず、200を返してしまうのでここで書き込み https://stackoverflow.com/questions/76319196/unit-testing-of-gins-context-redirect-works-for-get-response-code-but-fails-for

		assert.Equal(t, "/signup", w.Header().Get("Location"))
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
		values.Add("name", "test")
		values.Add("email", "test@example.com")
		values.Add("password", "test1234")
		c.Request = httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockCreateUserUsecase)
		mockUC.On("Invoke", mock.Anything).Return(errors.New("internal error"))

		createUser(c, mockUC, s)
		c.Writer.WriteHeaderNow() // POST だとヘッダの書き込みが行われず、200を返してしまうのでここで書き込み https://stackoverflow.com/questions/76319196/unit-testing-of-gins-context-redirect-works-for-get-response-code-but-fails-for

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockUC.AssertExpectations(t)
	})
}
