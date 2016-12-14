package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/server"
)

func setup(testUser *model.User, otherTestUser *model.User, s **gin.Engine, baseURL string) {
	testConfig := config.GetConfig("test")

	db := database.NewConnection("5432")
	redis := cache.NewRedisConnection("6379")

	gin.SetMode(gin.TestMode)
	db.Exec("TRUNCATE users, tweets CASCADE;") // Ugly, but lets keep it for convenience for now

	err := db.QueryRow("INSERT INTO users (username, email, password, name)"+
		"VALUES ($1, $2, $3, $4) RETURNING id, username, email, password, name",
		"user", "user@email.com", "password", "name").
		Scan(&testUser.ID, &testUser.Username, &testUser.Email, &testUser.Password, &testUser.Name)
	if err != nil {
		panic(fmt.Sprintf("Error inserting test user into database = %v", err))
	}

	err = db.QueryRow("INSERT INTO users (username, email, password, name)"+
		"VALUES ($1, $2, $3, $4) RETURNING id, username, email, password, name",
		"otheruser", "otheruser@email.com", "otherpassword", "othername").
		Scan(&otherTestUser.ID, &otherTestUser.Username, &otherTestUser.Email, &otherTestUser.Password, &otherTestUser.Name)
	if err != nil {
		panic(fmt.Sprintf("Error inserting other test user into database = %v", err))
	}

	*s = server.New(db, redis, testConfig)

	baseURL = "http://localhost:8080"
}

func createUser(name string, s *gin.Engine, url string, t *testing.T) *model.User {
	newUserForm := model.NewUserForm{
		Email:    name + "@email.com",
		Password: name + "password",
		Name:     name + "name",
		Username: name + "username",
	}
	data, _ := json.Marshal(newUserForm)

	buf := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", url+"/signup", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("error logging user int, status code: %v, expected: %v", w.Code, http.StatusOK)
	}

	var newUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &newUser)
	if err != nil {
		t.Error(err)
	}

	return &model.User{
		ID:       newUser.ID,
		Username: newUser.Username,
		Name:     newUser.Name,
		Email:    newUserForm.Email,
		Password: newUserForm.Password,
	}
}

func loginUser(user *model.User, s *gin.Engine, url string, t *testing.T) (string, string) {
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

	return loginResponse.AuthToken, loginResponse.RefreshToken
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

func likeTweet(tweetID int64, authToken string, s *gin.Engine, url string, t *testing.T) {
	req, _ := http.NewRequest("POST", url+"/tweets/"+strconv.FormatInt(int64(tweetID), 10)+"/like", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error liking tweet, status code: %v, expected: %v", w.Code, http.StatusOK)
	}
}

func unlikeTweet(tweetID int64, authToken string, s *gin.Engine, url string, t *testing.T) {
	req, _ := http.NewRequest("POST", url+"/tweets/"+strconv.FormatInt(int64(tweetID), 10)+"/unlike", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error unliking tweet, status code: %v, expected: %v", w.Code, http.StatusOK)
	}
}

func followUser(userID int64, authToken string, s *gin.Engine, url string, t *testing.T) {
	req, _ := http.NewRequest("POST", url+"/users/"+strconv.FormatInt(int64(userID), 10)+"/follow", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error following user, status code: %v, expected: %v", w.Code, http.StatusOK)
	}
}

func unfollowUser(userID int64, authToken string, s *gin.Engine, url string, t *testing.T) {
	req, _ := http.NewRequest("POST", url+"/users/"+strconv.FormatInt(int64(userID), 10)+"/unfollow", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error unfollowing user, status code: %v, expected: %v", w.Code, http.StatusOK)
	}
}
