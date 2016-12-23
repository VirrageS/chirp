package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/server"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration")
}

var _ = Describe("ServerTest", func() {
	var (
		router *gin.Engine
		db     *sql.DB

		ala             *model.User
		bob             *model.User
		toor            *model.User
		ernest          *model.User
		alaToken        string
		alaRefreshToken string
		bobToken        string
		alaPublic       *model.PublicUser
		bobPublic       *model.PublicUser
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)

		db = database.NewConnection("5433")
		testConfig := config.GetConfig("test")
		cache := cache.NewDummyCache()
		router = server.New(db, cache, testConfig)

		// create users
		ala = createUser(router, "ala")
		bob = createUser(router, "bob")
		alaToken, alaRefreshToken = loginUser(router, ala)
		bobToken, _ = loginUser(router, bob)

		alaPublic = retrieveUser(router, ala.ID, alaToken)
		bobPublic = retrieveUser(router, bob.ID, bobToken)

		// create additional users
		toor = createUser(router, "toor")
		ernest = createUser(router, "ernest")
	})

	AfterEach(func() {
		// HACK: this is hack since TRUNCATE can execute up to 1s... whereas this ~5ms
		db.Exec(`
			DELETE FROM users;
			DELETE FROM tweets;
			DELETE FROM follows;
			DELETE FROM likes;
			DELETE FROM retweets;
		`)
	})

	Describe("Create new user", func() {
		var (
			newUserForm *model.NewUserForm
		)

		BeforeEach(func() {
			newUserForm = &model.NewUserForm{
				Username: "anotherUser",
				Password: "anotherPassword",
				Email:    "another@email.com",
				Name:     "anotherName",
			}
		})

		It("should create user and populate fields correctly", func() {
			req := request("POST", "/signup", body(newUserForm)).json().build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var newUser model.PublicUser
			err := json.Unmarshal(w.Body.Bytes(), &newUser)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(newUser.Username).To(Equal(newUserForm.Username))
			Expect(newUser.Name).To(Equal(newUserForm.Name))
			Expect(newUser.AvatarUrl).To(Equal(""))
			Expect(newUser.Following).To(Equal(false))
			Expect(newUser.FollowerCount).To(BeEquivalentTo(0))
		})

		It("should not create user when other user exists with same username", func() {
			newUserForm.Username = ala.Username

			req := request("POST", "/signup", body(newUserForm)).json().build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusConflict))
			Expect(w.Body.Len()).NotTo(BeEquivalentTo(0))
		})

		It("should not create user when other user exists with same email", func() {
			newUserForm.Email = ala.Email

			req := request("POST", "/signup", body(newUserForm)).json().build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusConflict))
			Expect(w.Body.Len()).NotTo(BeEquivalentTo(0))
		})
	})

	Describe("Login user", func() {
		var (
			loggedUser    *model.PublicUser
			loginUserForm *model.LoginForm
		)

		BeforeEach(func() {
			loggedUser = alaPublic
			loginUserForm = &model.LoginForm{
				Email:    ala.Email,
				Password: ala.Password,
			}
		})

		It("should login user and return logged in user", func() {
			req := request("POST", "/login", body(loginUserForm)).json().build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var loginResponse model.LoginResponse
			err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(loginResponse.User).To(Equal(loggedUser))
			Expect(loginResponse.AuthToken).NotTo(BeEmpty())
			Expect(loginResponse.RefreshToken).NotTo(BeEmpty())
		})

		It("should not login user with wrong password", func() {
			loginUserForm.Password = "invalidpassword"

			req := request("POST", "/login", body(loginUserForm)).json().build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
			Expect(w.Body.Len()).NotTo(BeEquivalentTo(0))
		})

		It("should not login user with wrong email", func() {
			loginUserForm.Email = "invalid@email.com"

			req := request("POST", "/login", body(loginUserForm)).json().build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
			Expect(w.Body.Len()).NotTo(BeEquivalentTo(0))
		})
	})

	Describe("Follow user", func() {
		BeforeEach(func() {})

		It("should follow user and populate fields appropriately", func() {
			path := fmt.Sprintf("/users/%v/follow", toor.ID)
			req := request("POST", path, nil).authorize(alaToken).build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var actualUser model.PublicUser
			err := json.Unmarshal(w.Body.Bytes(), &actualUser)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(actualUser.FollowerCount).To(BeEquivalentTo(1))
			Expect(actualUser.Following).To(BeTrue())
		})

		It("should followed user match actual user", func() {
			actualUser := followUser(router, toor.ID, alaToken)
			expectedUser := retrieveUser(router, toor.ID, alaToken)
			Expect(actualUser).To(Equal(expectedUser))
		})

		It("should not update follow user nor unfollow when following user twice", func() {
			followUser(router, toor.ID, alaToken)
			followUser(router, toor.ID, alaToken)

			actualUser := retrieveUser(router, toor.ID, alaToken)
			Expect(actualUser.FollowerCount).To(BeEquivalentTo(1))
			Expect(actualUser.Following).To(BeTrue())
		})

		It("should update follow when two different user follow other user", func() {
			followUser(router, toor.ID, alaToken)
			followUser(router, toor.ID, bobToken)

			actualUser := retrieveUser(router, toor.ID, alaToken)
			Expect(actualUser.FollowerCount).To(BeEquivalentTo(2))
			Expect(actualUser.Following).To(BeTrue())
		})
	})

	Describe("Unfollow user", func() {
		BeforeEach(func() {})

		It("should unfollow user which is followed and populate all data", func() {
			followUser(router, toor.ID, alaToken)

			path := fmt.Sprintf("/users/%v/unfollow", toor.ID)
			req := request("POST", path, nil).authorize(alaToken).build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var actualUser model.PublicUser
			err := json.Unmarshal(w.Body.Bytes(), &actualUser)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(actualUser.FollowerCount).To(BeEquivalentTo(0))
			Expect(actualUser.Following).To(BeFalse())
		})

		It("should unfollowed user match actual user", func() {
			followUser(router, toor.ID, alaToken)

			actualUser := unfollowUser(router, toor.ID, alaToken)
			expectedUser := retrieveUser(router, toor.ID, alaToken)
			Expect(actualUser).To(Equal(expectedUser))
		})

		It(`should not perform any operation (but should return user)
				when trying to unfollow not followed user`, func() {
			unfollowUser(router, toor.ID, alaToken)
			actualUser := retrieveUser(router, toor.ID, alaToken)
			Expect(actualUser.FollowerCount).To(BeEquivalentTo(0))
			Expect(actualUser.Following).To(BeFalse())
		})

		It(`should not perform any operation (but should return user)
				when trying unfollow user followed by someone else`, func() {
			followUser(router, toor.ID, alaToken)
			unfollowUser(router, toor.ID, bobToken)

			alaActualUser := retrieveUser(router, toor.ID, alaToken)
			Expect(alaActualUser.FollowerCount).To(BeEquivalentTo(1))
			Expect(alaActualUser.Following).To(BeTrue())

			bobActualUser := retrieveUser(router, toor.ID, bobToken)
			Expect(bobActualUser.FollowerCount).To(BeEquivalentTo(1))
			Expect(bobActualUser.Following).To(BeFalse())
		})
	})

	Describe("Get followers", func() {
		BeforeEach(func() {})

		It("should get followers of followed user", func() {
			expectedFollowers := []*model.PublicUser{
				alaPublic,
			}

			followUser(router, toor.ID, alaToken)

			actualFollowers := retrieveFollowers(router, toor.ID, alaToken)
			Expect(*actualFollowers).To(Equal(expectedFollowers))
		})

		It("should get followers of user followed by multiple users", func() {
			expectedFollowers := []*model.PublicUser{
				alaPublic,
				bobPublic,
			}

			followUser(router, toor.ID, alaToken)
			followUser(router, toor.ID, bobToken)

			actualFollowers := retrieveFollowers(router, toor.ID, alaToken)
			Expect(*actualFollowers).To(Equal(expectedFollowers))
		})

		It("should get followers of not followed user", func() {
			actualFollowers := retrieveFollowers(router, toor.ID, alaToken)
			Expect(*actualFollowers).To(BeEmpty())
		})

		It("should get followers of user followed by someone else", func() {
			expectedFollowers := []*model.PublicUser{bobPublic}

			followUser(router, toor.ID, bobToken)

			actualFollowers := retrieveFollowers(router, toor.ID, alaToken)
			Expect(*actualFollowers).To(ConsistOf(expectedFollowers))
		})
	})

	Describe("Get followees", func() {
		BeforeEach(func() {})

		It("should get followees", func() {
			expectedFollowees := []*model.PublicUser{
				publicUser(*toor).followerCount(1).following(true).build(),
				publicUser(*ernest).followerCount(1).following(true).build(),
			}

			followUser(router, toor.ID, alaToken)
			followUser(router, ernest.ID, alaToken)

			actualFollowees := retrieveFollowees(router, ala.ID, alaToken)
			Expect(*actualFollowees).To(ConsistOf(expectedFollowees))
		})

		It("should get empty followees when user is not following anyone", func() {
			actualFollowers := retrieveFollowees(router, ala.ID, alaToken)
			Expect(*actualFollowers).To(BeEmpty())
		})

		It("should get only current user followees", func() {
			expectedFollowees := []*model.PublicUser{
				publicUser(*toor).followerCount(1).following(true).build(),
			}

			followUser(router, toor.ID, alaToken)
			followUser(router, ernest.ID, bobToken)

			actualFollowers := retrieveFollowees(router, ala.ID, alaToken)
			Expect(*actualFollowers).To(ConsistOf(expectedFollowees))
		})
	})

	Describe("Create and get tweet", func() {
		BeforeEach(func() {})

		It("should create tweet", func() {
			newTweet := &model.NewTweet{Content: "new tweet"}
			req := request("POST", "/tweets", body(newTweet)).json().authorize(alaToken).build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var actualTweet model.Tweet
			err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(actualTweet.LikeCount).To(BeEquivalentTo(0))
			Expect(actualTweet.RetweetCount).To(BeEquivalentTo(0))
			Expect(actualTweet.Content).To(Equal("new tweet"))
			Expect(actualTweet.Liked).To(Equal(false))
			Expect(actualTweet.Retweeted).To(Equal(false))
			Expect(actualTweet.Author).To(Equal(alaPublic))
		})

		It("should get tweet after creating", func() {
			expectedTweet := createTweet(router, "new tweet", alaToken)
			actualTweet := retrieveTweet(router, expectedTweet.ID, alaToken)

			Expect(actualTweet).To(Equal(expectedTweet))
		})
	})

	Describe("Delete tweet", func() {
		BeforeEach(func() {})

		It("should delete existing tweet", func() {
			createdTweet := createTweet(router, "new tweet", alaToken)
			deleteTweet(router, createdTweet.ID, alaToken)

			path := fmt.Sprintf("/tweets/%v", createdTweet.ID)
			req := request("GET", path, nil).authorize(alaToken).build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("should not perform any unexpected actions trying to delete not existing tweet", func() {
			createdTweet := createTweet(router, "new tweet", alaToken)
			deleteTweet(router, createdTweet.ID, alaToken)
		})

		It("should not allow to delete tweet created by someone else", func() {
			createdTweet := createTweet(router, "new tweet", bobToken)

			path := fmt.Sprintf("/tweets/%v", createdTweet.ID)
			req := request("DELETE", path, nil).authorize(alaToken).build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusForbidden))
		})
	})

	Describe("Get tweets", func() {
		BeforeEach(func() {})

		It("should get tweets after creation", func() {
			expectedTweets := []*model.Tweet{
				createTweet(router, "new tweet1", alaToken),
				createTweet(router, "new tweet2", alaToken),
			}

			actualTweets := retrieveTweets(router, alaToken)
			Expect(*actualTweets).To(ConsistOf(expectedTweets))
		})

		It(`should only get tweets created by specified user
				when providing user in query parameter`, func() {
			alaExpectedTweets := []*model.Tweet{
				createTweet(router, "tweet1", alaToken),
				createTweet(router, "tweet2", alaToken),
			}

			bobExpectedTweets := []*model.Tweet{
				createTweet(router, "something different", bobToken),
			}

			alaActualTweets := retrieveUserTweets(router, alaToken, ala.ID)
			Expect(*alaActualTweets).To(ConsistOf(alaExpectedTweets))

			bobActualTweets := retrieveUserTweets(router, bobToken, bob.ID)
			Expect(*bobActualTweets).To(ConsistOf(bobExpectedTweets))
		})
	})

	Describe("Get home feed", func() {
		BeforeEach(func() {})

		It("should get tweets created by user", func() {
			alaExpectedTweets := []*model.Tweet{
				createTweet(router, "tweet1", alaToken),
				createTweet(router, "tweet2", alaToken),
			}

			bobExpectedTweets := []*model.Tweet{
				createTweet(router, "something different", bobToken),
			}

			alaActualTweets := retrieveHomeFeed(router, alaToken)
			Expect(*alaActualTweets).To(ConsistOf(alaExpectedTweets))

			bobActualTweets := retrieveHomeFeed(router, bobToken)
			Expect(*bobActualTweets).To(ConsistOf(bobExpectedTweets))
		})
	})

	Describe("Like tweet", func() {
		var (
			alaTweet *model.Tweet
			bobTweet *model.Tweet
		)

		BeforeEach(func() {
			alaTweet = createTweet(router, "new ala tweet", alaToken)
			bobTweet = createTweet(router, "new bob tweet", bobToken)
		})

		It("should like tweet and return new liked tweet with populated data", func() {
			actualTweet := likeTweet(router, alaTweet.ID, alaToken)
			Expect(actualTweet.LikeCount).To(BeEquivalentTo(1))
			Expect(actualTweet.Liked).To(Equal(true))
		})

		It("should liked tweet match actual tweet", func() {
			actualTweet := likeTweet(router, alaTweet.ID, alaToken)
			expectedTweet := retrieveTweet(router, alaTweet.ID, alaToken)
			Expect(actualTweet).To(Equal(expectedTweet))
		})

		It("should tweet be liked only once after consecutive likes", func() {
			likeTweet(router, alaTweet.ID, alaToken)
			likeTweet(router, alaTweet.ID, alaToken)

			actualTweet := retrieveTweet(router, alaTweet.ID, alaToken)
			Expect(actualTweet.LikeCount).To(BeEquivalentTo(1))
			Expect(actualTweet.Liked).To(Equal(true))
		})

		It("should return tweet liked twice after multiple (two) likes", func() {
			likeTweet(router, alaTweet.ID, alaToken)
			likeTweet(router, alaTweet.ID, bobToken)

			actualTweet := retrieveTweet(router, alaTweet.ID, alaToken)
			Expect(actualTweet.LikeCount).To(BeEquivalentTo(2))
		})

		It("should return tweet with false `liked` field when tweet was not liked by user", func() {
			likeTweet(router, bobTweet.ID, bobToken)

			actualTweet := retrieveTweet(router, bobTweet.ID, alaToken)
			Expect(actualTweet.Liked).To(Equal(false))
		})
	})

	Describe("Unlike tweet", func() {
		var (
			alaTweet *model.Tweet
			bobTweet *model.Tweet
		)

		BeforeEach(func() {
			alaTweet = createTweet(router, "new ala tweet", alaToken)
			bobTweet = createTweet(router, "new bob tweet", bobToken)
		})

		It("should unlike tweet and return unliked tweet with fresh data", func() {
			likeTweet(router, alaTweet.ID, alaToken)

			path := fmt.Sprintf("/tweets/%v/unlike", alaTweet.ID)
			req := request("POST", path, nil).authorize(alaToken).build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var actualTweet model.Tweet
			err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(actualTweet.LikeCount).To(BeEquivalentTo(0))
			Expect(actualTweet.Liked).To(Equal(false))
		})

		It("should new unliked tweet match actual tweet", func() {
			likeTweet(router, alaTweet.ID, alaToken)
			actualTweet := unlikeTweet(router, alaTweet.ID, alaToken)
			expectedTweet := retrieveTweet(router, alaTweet.ID, alaToken)
			Expect(actualTweet).To(Equal(expectedTweet))
		})

		It("should not perform any unexpected actions when trying to unlike not liked tweet", func() {
			actualTweet := unlikeTweet(router, alaTweet.ID, alaToken)
			expectedTweet := retrieveTweet(router, alaTweet.ID, alaToken)
			Expect(actualTweet).To(Equal(expectedTweet))
		})

		It("should not unlike tweet which is liked by someone else", func() {
			likeTweet(router, alaTweet.ID, alaToken)
			unlikeTweet(router, alaTweet.ID, bobToken)

			alaActualTweet := retrieveTweet(router, alaTweet.ID, alaToken)
			Expect(alaActualTweet.LikeCount).To(BeEquivalentTo(1))
			Expect(alaActualTweet.Liked).To(Equal(true))

			bobActualTweet := retrieveTweet(router, alaTweet.ID, bobToken)
			Expect(bobActualTweet.LikeCount).To(BeEquivalentTo(1))
			Expect(bobActualTweet.Liked).To(Equal(false))
		})
	})

	Describe("Refresh auth token", func() {
		It("should refresh auth token", func() {
			refreshTokenRequest := &model.RefreshAuthTokenRequest{
				UserID:       ala.ID,
				RefreshToken: alaRefreshToken,
			}

			req := request("POST", "/token", body(refreshTokenRequest)).json().build()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))

			var refreshResponse model.RefreshAuthTokenResponse
			err := json.Unmarshal(w.Body.Bytes(), &refreshResponse)
			Expect(err).NotTo(HaveOccurred())

			newAuthToken := refreshResponse.AuthToken
			Expect(newAuthToken).NotTo(BeEmpty())

			// test creating tweet with new auth
			createdTweet := createTweet(router, "new tweet", newAuthToken)
			Expect(createdTweet.Author).To(Equal(alaPublic))
		})
	})
})
