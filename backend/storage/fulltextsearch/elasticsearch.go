package fulltextsearch

import (
	"context"
	"encoding/json"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/VirrageS/chirp/backend/config"
	"gopkg.in/olivere/elastic.v5"
)

const (
	indexName = "fts"

	tweetType         = "tweet"
	tweetContentField = "content"

	userType          = "user"
	userNameField     = "name"
	userUsernameField = "username"
)

type elasticsearchClient struct {
	*elastic.Client
}

// NewElasticsearchSearch connects to Elasticsearch and returns new instance of
// struct which implements Searcher functions.
func NewElasticsearchSearch(config config.ElasticsearchConfigProvider) Searcher {
	username := config.GetUsername()
	password := config.GetPassword()
	url := fmt.Sprintf("http://%v:%v", config.GetHost(), config.GetPort())

	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetBasicAuth(username, password),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.WithError(err).Error("Error connecting to elasticsearch.")
		return nil
	}

	return &elasticsearchClient{client}
}

func (e *elasticsearchClient) GetTweetsIDs(querystring string) ([]int64, error) {
	return e.getIDsFromIndex(querystring, tweetType, tweetContentField)
}

func (e *elasticsearchClient) GetUsersIDs(querystring string) ([]int64, error) {
	return e.getIDsFromIndex(querystring, userType, userUsernameField, userNameField)
}

func (e *elasticsearchClient) getIDsFromIndex(querystring, typeName string, fields ...string) ([]int64, error) {
	// Creates a MatchQuery with "and" operator - a query that will require
	// each word in `querystring` to be matched in one of the `fields`.
	query := elastic.NewMultiMatchQuery(querystring, fields...).Operator("and")

	searchResult, err := e.Search().
		Index(indexName).
		Type(typeName).
		Query(query).
		Do(context.Background())

	if err != nil {
		log.WithError(err).Error("Error querying elasticsearch in getIDsFromIndex.")
		return nil, err
	}

	ids := make([]int64, searchResult.Hits.TotalHits)

	for i, hit := range searchResult.Hits.Hits {
		var id struct {
			ID int64 `json:"id" binding:"required"`
		}

		err := json.Unmarshal(*hit.Source, &id)
		if err != nil {
			log.WithError(err).Error("Error umarshalling elasticsearch response in getIDsFromIndex.")
			return nil, err
		}

		ids[i] = id.ID
	}

	// Results from elasticsearch are sorted by score by default.
	// We do not need to sort them ourselves.
	return ids, nil
}
