package database

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/VirrageS/chirp/backend/database/model"
)

// Struct that implements TweetDataAccessor using sql (postgres) database
type TweetDB struct {
	*sql.DB
}

// Constructs TweetDB that uses a given sql.DB connection
func NewTweetDB(databaseConnection *sql.DB) *TweetDB {
	return &TweetDB{databaseConnection}
}

func (db *TweetDB) GetTweets() ([]model.TweetWithAuthor, error) {
	tweets, err := db.getTweets()
	if err != nil {
		return nil, DatabaseError
	}

	return tweets, nil
}

func (db *TweetDB) GetTweetsOfUserWithID(userID int64) ([]model.TweetWithAuthor, error) {
	tweets, err := db.getTweetsOfUserWithID(userID)
	if err != nil {
		return nil, DatabaseError
	}

	return tweets, nil
}

func (db *TweetDB) GetTweet(tweetID int64) (model.TweetWithAuthor, error) {
	tweet, err := db.getTweetUsingQuery(
		"SELECT tweets.id, tweets.created_at, tweets.content, "+
			"users.id, users.username, users.last_login, users.name, users.avatar_url "+
			"FROM tweets JOIN users on tweets.author_id=users.id AND tweets.id=$1;", tweetID)
	if err == sql.ErrNoRows {
		return model.TweetWithAuthor{}, NoResults
	}
	if err != nil {
		return model.TweetWithAuthor{}, DatabaseError
	}

	return tweet, nil
}

func (db *TweetDB) InsertTweet(tweet model.Tweet) (model.TweetWithAuthor, error) {
	tweetID, err := db.insertTweetToDatabase(tweet)
	if err != nil {
		return model.TweetWithAuthor{}, DatabaseError
	}

	// TODO: this is probably super ugly, fix it
	newTweet, err := db.getTweetUsingQuery(
		"SELECT tweets.id, tweets.created_at, tweets.content, "+
			"users.id, users.username, users.last_login, users.name, users.avatar_url "+
			"FROM tweets JOIN users on tweets.author_id=users.id AND tweets.id=$1;", tweetID)
	if err != nil {
		return model.TweetWithAuthor{}, DatabaseError
	}

	return newTweet, nil
}

func (db *TweetDB) DeleteTweet(tweetID int64) error {
	err := db.deleteTweetWithID(tweetID)
	if err != nil {
		return DatabaseError
	}

	return nil
}

// TODO: Maybe it should also fetch tweet's User and embed it inside the returned object
func (db *TweetDB) getTweetUsingQuery(query string, args ...interface{}) (model.TweetWithAuthor, error) {
	row := db.QueryRow(query, args...)

	var tweet model.TweetWithAuthor
	var author model.PublicUser

	err := row.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
		&author.ID, &author.Username, &author.LastLogin, &author.Name, &author.AvatarUrl)
	if err != nil && err != sql.ErrNoRows {
		log.WithField("query", query).WithError(err).Error("getTweetUsingQuery database error.")
	}
	tweet.Author = &author

	return tweet, err
}

func (db *TweetDB) insertTweetToDatabase(tweet model.Tweet) (int64, error) {
	query, err := db.Prepare("INSERT INTO tweets (author_id, created_at, content) " +
		"VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		log.WithError(err).Error("insertTweetToDatabase query prepare error.")
		return 0, err
	}
	defer query.Close()

	var newID int64

	err = query.QueryRow(tweet.AuthorID, tweet.CreatedAt, tweet.Content).Scan(&newID)
	if err != nil {
		log.WithError(err).Error("insertTweetToDatabase query execute error.")
		return 0, err
	}

	return newID, nil
}

func (db *TweetDB) deleteTweetWithID(tweetID int64) error {
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

func (db *TweetDB) getTweets() ([]model.TweetWithAuthor, error) {
	rows, err := db.Query(
		"SELECT tweets.id, tweets.created_at, tweets.content, " +
			"users.id, users.username, users.last_login, users.name, users.avatar_url " +
			"FROM tweets JOIN users on tweets.author_id=users.id;")
	if err != nil {
		log.WithError(err).Error("getTweets query error.")
	}

	var tweets []model.TweetWithAuthor

	defer rows.Close()
	for rows.Next() {
		var tweet model.TweetWithAuthor
		var author model.PublicUser

		err := rows.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
			&author.ID, &author.Username, &author.LastLogin, &author.Name, &author.AvatarUrl)
		if err != nil {
			log.WithError(err).Error("getTweets row scan error.")
			return nil, err
		}
		tweet.Author = &author

		tweets = append(tweets, tweet)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("getTweets rows iteration error.")
		return nil, err
	}

	return tweets, nil
}

// TODO: almost the same as getTweets()...
func (db *TweetDB) getTweetsOfUserWithID(userID int64) ([]model.TweetWithAuthor, error) {
	rows, err := db.Query("SELECT tweets.id, tweets.created_at, tweets.content, "+
		"users.id, users.username, users.last_login, users.name, users.avatar_url "+
		"FROM tweets JOIN users on tweets.author_id=users.id AND users.id=$1;", userID)
	if err != nil {
		log.WithError(err).Error("getTweetsOfUserWithID query error.")
	}

	var tweets []model.TweetWithAuthor

	defer rows.Close()
	for rows.Next() {
		var tweet model.TweetWithAuthor
		var author model.PublicUser

		err := rows.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
			&author.ID, &author.Username, &author.LastLogin, &author.Name, &author.AvatarUrl)
		if err != nil {
			log.WithError(err).Error("getTweetsOfUserWithID row scan error.")
			return nil, err
		}
		tweet.Author = &author

		tweets = append(tweets, tweet)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("getTweetsOfUserWithID rows iteration error.")
		return nil, err
	}

	return tweets, nil
}
