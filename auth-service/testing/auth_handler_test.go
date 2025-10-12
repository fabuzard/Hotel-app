package testing

import (
	"auth-service/dto"
	"auth-service/handler"
	"auth-service/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ---- Mock Service ----
type dummyAuthService struct{}

func (d *dummyAuthService) RegisterUser(req dto.RegisterRequest) (model.Users, error) {
	if req.Email == "exists@test.com" {
		return model.Users{}, errors.New("email already registered")
	}
	if req.Password == "" {
		return model.Users{}, errors.New("password cannot be empty")
	}
	return model.Users{ID: 1, Email: req.Email, FullName: req.FullName}, nil
}

func (d *dummyAuthService) LoginRequest(req dto.LoginRequest) (string, error) {
	if req.Email == "wrong@test.com" {
		return "", errors.New("invalid email or password")
	}
	if req.Password == "" {
		return "", errors.New("password required")
	}
	return "dummy-jwt-token", nil
}

// ---- Test: Register (Success) ----
func TestRegister_Success(t *testing.T) {
	e := echo.New()
	h := handler.NewAuthHandler(&dummyAuthService{})

	body := `{"email":"new@test.com","password":"123","full_name":"Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.Register(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "new@test.com")
}

// ---- Test: Register (Already Exists) ----
func TestRegister_AlreadyExists(t *testing.T) {
	e := echo.New()
	h := handler.NewAuthHandler(&dummyAuthService{})

	body := `{"email":"exists@test.com","password":"123","full_name":"Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.Register(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "already registered")
}

// ---- Test: Register (Missing Password) ----
func TestRegister_MissingPassword(t *testing.T) {
	e := echo.New()
	h := handler.NewAuthHandler(&dummyAuthService{})

	body := `{"email":"new@test.com","full_name":"Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.Register(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "password")
}

// ---- Test: Login (Success) ----
func TestLogin_Success(t *testing.T) {
	e := echo.New()
	h := handler.NewAuthHandler(&dummyAuthService{})

	body := `{"email":"new@test.com","password":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.Login(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "dummy-jwt-token")
}

// ---- Test: Login (Wrong Email) ----
func TestLogin_WrongEmail(t *testing.T) {
	e := echo.New()
	h := handler.NewAuthHandler(&dummyAuthService{})

	body := `{"email":"wrong@test.com","password":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.Login(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid email or password")
}

// ---- Test: Login (Missing Password) ----
func TestLogin_MissingPassword(t *testing.T) {
	e := echo.New()
	h := handler.NewAuthHandler(&dummyAuthService{})

	body := `{"email":"new@test.com"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.Login(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "password")
}
