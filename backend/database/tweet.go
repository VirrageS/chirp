package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
)

// Struct that implements TweetDataAccessor using sql (postgres) database
type TweetDB struct {
	*sql.DB
}

// Constructs TweetDB that uses a given sql.DB connection
func NewTweetDB(databaseConnection *sql.DB) *TweetDB {
	return &TweetDB{databaseConnection}
}

func (db *TweetDB) GetTweets() ([]*model.Tweet, error) {
	tweets, err := db.getTweets()
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweets, nil
}

func (db *TweetDB) GetTweetsOfUserWithID(userID int64) ([]*model.Tweet, error) {
	tweets, err := db.getTweetsOfUserWithID(userID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweets, nil
}

func (db *TweetDB) GetTweet(tweetID int64) (*model.Tweet, error) {
	tweet, err := db.getTweetUsingQuery(
		"SELECT tweets.id, tweets.created_at, tweets.content, "+
			"users.id, users.username, users.name, users.avatar_url "+
			"FROM tweets JOIN users on tweets.author_id=users.id AND tweets.id=$1;", tweetID)
	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweet, nil
}

func (db *TweetDB) InsertTweet(tweet *model.NewTweet) (*model.Tweet, error) {
	tweetID, err := db.insertTweetToDatabase(tweet)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	// TODO: this is probably super ugly. Maybe fetch user only?
	// Probably could just fetch user from cache
	newTweet, err := db.getTweetUsingQuery(
		"SELECT tweets.id, tweets.created_at, tweets.content, "+
			"users.id, users.username, users.name, users.avatar_url "+
			"FROM tweets JOIN users on tweets.author_id=users.id AND tweets.id=$1;", tweetID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return newTweet, nil
}

func (db *TweetDB) DeleteTweet(tweetID int64) error {
	err := db.deleteTweetWithID(tweetID)
	if err != nil {
		return errors.UnexpectedError
	}

	return nil
}

func (db *TweetDB) LikeTweet(tweetID, userID int64) error {
	err := db.likeTweet(tweetID, userID)
	if err != nil {
		return errors.UnexpectedError
	}

	return nil
}

// TODO: Maybe it should also fetch tweet's User and embed it inside the returned object
func (db *TweetDB) getTweetUsingQuery(query string, args ...interface{}) (*model.Tweet, error) {
	row := db.QueryRow(query, args...)

	var tweet model.Tweet
	var author model.PublicUser

	err := row.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
		&author.ID, &author.Username, &author.Name, &author.AvatarUrl)
	if err != nil && err != sql.ErrNoRows {
		log.WithField("query", query).WithError(err).Error("getTweetUsingQuery database error.")
		return nil, err
	}
	tweet.Author = &author

	return &tweet, err
}

// TODO: maybe return whole Tweet struct instead of just ID
func (db *TweetDB) insertTweetToDatabase(tweet *model.NewTweet) (int64, error) {
	query, err := db.Prepare("INSERT INTO tweets (author_id, content) " +
		"VALUES ($1, $2) RETURNING id")
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

func (db *TweetDB) getTweets() ([]*model.Tweet, error) {
	rows, err := db.Query(`
		SELECT tweets.id, tweets.created_at, tweets.content,
				users.id, users.username, users.name, users.avatar_url
		FROM tweets JOIN users on tweets.author_id=users.id
		ORDER BY tweets.created_at DESC;
	`)
	if err != nil {
		log.WithError(err).Error("getTweets query error.")
		return nil, err
	}

	var tweets []*model.Tweet

	defer rows.Close()
	for rows.Next() {
		var tweet model.Tweet
		var author model.PublicUser

		err := rows.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
			&author.ID, &author.Username, &author.Name, &author.AvatarUrl)
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
func (db *TweetDB) getTweetsOfUserWithID(userID int64) ([]*model.Tweet, error) {
	rows, err := db.Query(
		`SELECT tweets.id, tweets.created_at, tweets.content,
				users.id, users.username, users.name, users.avatar_url
		FROM tweets JOIN users on tweets.author_id=users.id AND users.id=$1
		ORDER BY tweets.created_at DESC;`,
		userID,
	)
	if err != nil {
		log.WithError(err).Error("getTweetsOfUserWithID query error.")
		return nil, err
	}

	var tweets []*model.Tweet

	defer rows.Close()
	for rows.Next() {
		var tweet model.Tweet
		var author model.PublicUser

		err := rows.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content,
			&author.ID, &author.Username, &author.Name, &author.AvatarUrl)
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

func (db *TweetDB) likeTweet(tweetID, userID int64) error {
	query, err := db.Prepare("INSERT INTO likes (tweet_id, user_id) " +
		"VALUES ($1, $2) " +
		"ON CONFLICT (tweet_id, user_id) DO NOTHING;")

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

	log.Error("EXECUTED QUERY")

	return nil
}
