package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
	"strconv"
	"testing"
)

func createUser(email string, password string, s *gin.Engine, url string) (int64, error) {
	loginForm := model.LoginForm{
		Email:    email,
		Password: password,
	}

	data, err := json.Marshal(loginForm)
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

func loginUser(user *model.User, s *gin.Engine, url string, t *testing.T) string {
	loginData := &model.LoginForm{
		Email:    user.Email,
		Password: user.Password,
	}

	data, _ := json.Marshal(loginData)

	buf := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", url+"/login", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error logging user int, status code: %v, expected: %v", w.Code, http.StatusOK)
	}

	var loginResponse model.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	if err != nil {
		t.Error(err)
	}

	return loginResponse.AuthToken
}

func createTweet(content string, authToken string, s *gin.Engine, url string, t *testing.T) *model.Tweet {
	newTweet1 := &model.NewTweet{
		Content: content,
	}
	data, _ := json.Marshal(newTweet1)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", url+"/tweets", buf)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("error creating tweet, status code: %v,  expected: %v", w.Code, http.StatusCreated)
	}

	var createdTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &createdTweet)
	if err != nil {
		t.Error(err)
	}

	return &createdTweet
}

func deleteTweet(tweetID int64, authToken string, s *gin.Engine, url string, t *testing.T) {
	reqDELETE, _ := http.NewRequest("DELETE", url+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqDELETE.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqDELETE)

	if w.Code != http.StatusNoContent {
		t.Errorf("error deleting tweet, status code: %v, expected: %v", w.Code, http.StatusNoContent)
	}
}
