package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/server"
	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080"

var s *gin.Engine

func TestMain(m *testing.M) {
	db := database.NewConnection("5555")
	db.Exec("TRUNCATE users, tweets CASCADE;") // Ugly, but lets keep it for convenience for now
	gin.SetMode(gin.TestMode)
	s = server.New(db)
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	newUser := &model.NewUserForm{
		Username: "user",
		Password: "password",
		Email:    "email@email.com",
		Name:     "name",
	}
	data, _ := json.Marshal(newUser)

	buf := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", baseURL+"/signup", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.Nil(t, err)

	// TODO: add ID check once database is cleaned after every run
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, actualUser.Username, newUser.Username)
	assert.Equal(t, actualUser.Name, newUser.Name)
	assert.Equal(t, actualUser.AvatarUrl, "")
	assert.Equal(t, actualUser.Following, false)
}

func TestLoginUser(t *testing.T) {
	newUser := &model.NewUserForm{
		Username: "user1",
		Password: "password1",
		Email:    "email1@email.com",
		Name:     "name1",
	}
	userID, err := createUser(newUser, s, baseURL)
	assert.Nil(t, err)

	loginData := &model.LoginForm{
		Email:    "email1@email.com",
		Password: "password1",
	}
	data, _ := json.Marshal(loginData)

	buf := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", baseURL+"/login", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var loginResponse model.LoginResponse
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.Nil(t, err)

	actualUser := loginResponse.User

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, actualUser.ID, userID)
	assert.Equal(t, actualUser.Username, newUser.Username)
	assert.Equal(t, actualUser.Name, newUser.Name)
	assert.Equal(t, actualUser.AvatarUrl, "")
	assert.Equal(t, actualUser.Following, false)
}
