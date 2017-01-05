package fulltextsearch

import (
	"context"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

const tweetIndex = "tweets"
const tweetType = "tweet"
const tweetContentField = "content"

const userIndex = "users"
const userType = "user"
const userNameField = "name"
const userUsernameField = "username"

type ElasticsearchClient struct {
	*elastic.Client
}

func NewElasticsearch() *ElasticsearchClient {
	client, err := elastic.NewClient()
	if err != nil {
		log.WithError(err).Fatal("Error connecting to elasticsearch.")
	}

	return &ElasticsearchClient{client}
}

func (e *ElasticsearchClient) GetTweetsIDs(querystring string) ([]int64, error) {
	return e.getIDsFromIndex(querystring, tweetIndex, tweetType, tweetContentField)
}

func (e *ElasticsearchClient) GetUsersIDs(querystring string) ([]int64, error) {
	return e.getIDsFromIndex(querystring, userIndex, userType, userUsernameField, userNameField)
}

func (e *ElasticsearchClient) getIDsFromIndex(querystring, indexName, typeName string, fields ...string) ([]int64, error) {
	// Create a MatchQuery with "and" operator - a query that will require each word in `querystring` to be matched
	// in one of the `fields`.
	query := elastic.NewMultiMatchQuery(querystring, fields...).Operator("and")

	searchResult, err := e.Search().
		Index(indexName).
		Type(typeName).
		Query(query).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	totalHits := searchResult.Hits.TotalHits
	IDs := make([]int64, totalHits)

	for i, hit := range searchResult.Hits.Hits {
		var ID idStruct

		err := json.Unmarshal(*hit.Source, &ID)
		if err != nil {
			return nil, err
		}

		IDs[i] = ID.ID
	}

	// results from elasticsearch are sorted by score by default, so we don't need to sort ourselves
	return IDs, nil
}
