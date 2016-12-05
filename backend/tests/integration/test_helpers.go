package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
)

func createUser(user *model.NewUserForm, s *gin.Engine, url string) (int64, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return 0, err
	}

	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", url+"/signup", buf)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		return 0, errors.New("")
	}

	var actualUser model.PublicUser
	err = json.Unmarshal(w.Body.Bytes(), &actualUser)
	if err != nil {
		return 0, err
	}

	return actualUser.ID, nil
}

func loginUser(loginData *model.LoginForm, s *gin.Engine, url string) (string, error) {
	data, _ := json.Marshal(loginData)

	buf := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", url+"/login", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return "", errors.New("")
	}

	var loginResponse model.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	if err != nil {
		return "", err
	}

	return loginResponse.AuthToken, nil
}
