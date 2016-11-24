package converters

import (
	"time"

	APIModel "github.com/VirrageS/chirp/backend/api/model"
	databaseModel "github.com/VirrageS/chirp/backend/database/model"
)

type TweetModelConverter interface {
	ConvertDatabaseTweetToAPITweet(tweet *databaseModel.TweetWithAuthor) *APIModel.Tweet
	ConvertAPINewTweetToDatabaseTweet(tweet *APIModel.NewTweet) *databaseModel.Tweet
	ConvertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets []*databaseModel.TweetWithAuthor) []*APIModel.Tweet
}

type TweetConverter struct {
	userConverter UserModelConverter
}

func NewTweetConverter(userConverter UserModelConverter) TweetModelConverter {
	return &TweetConverter{userConverter}
}

func (converter *TweetConverter) ConvertDatabaseTweetToAPITweet(tweet *databaseModel.TweetWithAuthor) *APIModel.Tweet {
	tweetID := tweet.ID
	author := tweet.Author
	likes := tweet.Likes
	retweets := tweet.Retweets
	createdAt := tweet.CreatedAt
	content := tweet.Content

	APIauthor := converter.userConverter.ConvertDatabasePublicUserToAPI(author)

	APITweet := &APIModel.Tweet{
		ID:        tweetID,
		Author:    APIauthor,
		Likes:     likes,
		Retweets:  retweets,
		CreatedAt: createdAt,
		Content:   content,
		Liked:     false,
		Retweeted: false,
	}
	return APITweet
}

func (converter *TweetConverter) ConvertAPINewTweetToDatabaseTweet(tweet *APIModel.NewTweet) *databaseModel.Tweet {
	authorId := tweet.AuthorID
	content := tweet.Content

	return &databaseModel.Tweet{
		ID:        0,
		AuthorID:  authorId,
		Likes:     0,
		Retweets:  0,
		CreatedAt: time.Now(),
		Content:   content,
	}
}

func (converter *TweetConverter) ConvertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets []*databaseModel.TweetWithAuthor) []*APIModel.Tweet {
	APITweets := make([]*APIModel.Tweet, 0)

	for _, databaseTweet := range databaseTweets {
		APITweet := converter.ConvertDatabaseTweetToAPITweet(databaseTweet)
		APITweets = append(APITweets, APITweet)
	}

	return APITweets
}
