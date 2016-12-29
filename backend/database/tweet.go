package database

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/VirrageS/chirp/backend/model"
)

type TweetDAO interface {
	GetTweetUsingQuery(query string, args ...interface{}) (*model.Tweet, error)
	InsertTweetToDatabase(tweet *model.NewTweet) (int64, error)
	DeleteTweetWithID(tweetID int64) error
	GetTweets(requestingUserID int64) ([]*model.Tweet, error)
	GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error)
	LikeTweet(tweetID, userID int64) error
	UnlikeTweet(tweetID, userID int64) error
}

type tweetDB struct {
	*sql.DB
}

func NewTweetDAO(dbConnection *sql.DB) TweetDAO {
	return &tweetDB{dbConnection}
}

// TODO: Maybe it should also fetch tweet's User and embed it inside the returned object
func (db *tweetDB) GetTweetUsingQuery(query string, args ...interface{}) (*model.Tweet, error) {
	row := db.QueryRow(query, args...)

	var tweet model.Tweet
	var author model.PublicUser

	err := row.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
		&author.ID, &author.Username, &author.Name, &author.AvatarUrl, &tweet.LikeCount, &tweet.Liked)
	if err != nil && err != sql.ErrNoRows {
		log.WithFields(log.Fields{
			"query": query,
			"args":  args,
		}).WithError(err).Error("getTweetUsingQuery database error.")
		return nil, err
	}
	tweet.Author = &author

	return &tweet, err
}

// TODO: maybe return whole Tweet struct instead of just ID
func (db *tweetDB) InsertTweetToDatabase(tweet *model.NewTweet) (int64, error) {
	query, err := db.Prepare(`
		INSERT INTO tweets (author_id, content) VALUES ($1, $2) RETURNING id;
	`)
	if err != nil {
		log.WithError(err).Error("insertTweetToDatabase query prepare error.")
		return 0, err
	}
	defer query.Close()

	var newID int64

	err = query.QueryRow(tweet.AuthorID, tweet.Content).Scan(&newID)
	if err != nil {
		log.WithError(err).Error("insertTweetToDatabase query execute error.")
		return 0, err
	}

	return newID, nil
}

func (db *tweetDB) DeleteTweetWithID(tweetID int64) error {
	statement, err := db.Prepare("DELETE FROM tweets WHERE id=$1")
	if err != nil {
		log.WithError(err).Error("deleteTweetWithID query prepare error.")
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(tweetID)
	if err != nil {
		log.WithError(err).Error("deleteTweetWithID query execute error.")
		return err
	}

	return nil
}

func (db *tweetDB) GetTweets(requestingUserID int64) ([]*model.Tweet, error) {
	rows, err := db.Query(`
		SELECT tweets.id, tweets.created_at, tweets.content,
		 	users.id, users.username, users.name, users.avatar_url,
		 	COUNT(likes.tweet_id) as likes,
		 	SUM(case when likes.user_id=$1 then 1 else 0 end) > 0 as liked
		FROM tweets
			JOIN users ON tweets.author_id = users.id
			LEFT JOIN likes ON tweets.id = likes.tweet_id
		GROUP BY tweets.id, users.id
		ORDER BY tweets.created_at DESC;`,
		requestingUserID)
	if err != nil {
		log.WithError(err).Error("getTweets query error.")
		return nil, err
	}
	defer rows.Close()

	tweets := make([]*model.Tweet, 0)
	for rows.Next() {
		var tweet model.Tweet
		var author model.PublicUser

		err := rows.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
			&author.ID, &author.Username, &author.Name, &author.AvatarUrl, &tweet.LikeCount, &tweet.Liked)
		if err != nil {
			log.WithError(err).Error("getTweets row scan error.")
			return nil, err
		}
		tweet.Author = &author

		tweets = append(tweets, &tweet)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("getTweets rows iteration error.")
		return nil, err
	}

	return tweets, nil
}

// TODO: almost the same as getTweets()...
func (db *tweetDB) GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error) {
	rows, err := db.Query(`
		SELECT tweets.id, tweets.created_at, tweets.content,
		 	users.id, users.username, users.name, users.avatar_url,
		 	COUNT(likes.tweet_id) as likes,
		 	SUM(case when likes.user_id=$1 then 1 else 0 end) > 0 as liked
		FROM tweets
			JOIN users ON tweets.author_id = users.id AND users.id=$2
			LEFT JOIN likes ON tweets.id = likes.tweet_id
		GROUP BY tweets.id, users.id
		ORDER BY tweets.created_at DESC;`,
		userID, requestingUserID)
	if err != nil {
		log.WithError(err).Error("getTweetsOfUserWithID query error.")
		return nil, err
	}
	defer rows.Close()

	tweets := make([]*model.Tweet, 0)
	for rows.Next() {
		var tweet model.Tweet
		var author model.PublicUser

		err := rows.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
			&author.ID, &author.Username, &author.Name, &author.AvatarUrl, &tweet.LikeCount, &tweet.Liked)
		if err != nil {
			log.WithError(err).Error("getTweetsOfUserWithID row scan error.")
			return nil, err
		}
		tweet.Author = &author

		tweets = append(tweets, &tweet)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("getTweetsOfUserWithID rows iteration error.")
		return nil, err
	}

	return tweets, nil
}

func (db *tweetDB) LikeTweet(tweetID, userID int64) error {
	query, err := db.Prepare(`
		INSERT INTO likes (tweet_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (tweet_id, user_id) DO NOTHING;
		`)

	if err != nil {
		log.WithError(err).Error("likeTweet query prepare error")
		return err
	}
	defer query.Close()

	_, err = query.Exec(tweetID, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"tweetID": tweetID,
			"userID":  userID,
		}).WithError(err).Error("likeTweet query execute error.")
		return err
	}

	return nil
}

func (db *tweetDB) UnlikeTweet(tweetID, userID int64) error {
	query, err := db.Prepare(`
		DELETE FROM likes
		WHERE tweet_id=$1 AND user_id=$2;
		`)

	if err != nil {
		log.WithError(err).Error("deleteTweet query prepare error")
		return err
	}
	defer query.Close()

	_, err = query.Exec(tweetID, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"tweetID": tweetID,
			"userID":  userID,
		}).WithError(err).Error("deleteTweet query execute error.")
		return err
	}

	return nil
}
