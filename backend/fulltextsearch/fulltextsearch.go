package fulltextsearch

import (
	"context"
	"encoding/json"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/VirrageS/chirp/backend/config"
	"gopkg.in/olivere/elastic.v5"
)

const indexName = "fts"

const tweetType = "tweet"
const tweetContentField = "content"

const userType = "user"
const userNameField = "name"
const userUsernameField = "username"

type ElasticsearchClient struct {
	*elastic.Client
}

func NewElasticsearch(config config.ElasticsearchConfigProvider) *ElasticsearchClient {
	username := config.GetUsername()
	password := config.GetPassword()
	url := fmt.Sprintf("http://%v:%v", config.GetHost(), config.GetPort())

	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetBasicAuth(username, password),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.WithError(err).Fatal("Error connecting to elasticsearch.")
	}

	return &ElasticsearchClient{client}
}

func (e *ElasticsearchClient) GetTweetsIDs(querystring string) ([]int64, error) {
	return e.getIDsFromIndex(querystring, tweetType, tweetContentField)
}

func (e *ElasticsearchClient) GetUsersIDs(querystring string) ([]int64, error) {
	return e.getIDsFromIndex(querystring, userType, userUsernameField, userNameField)
}

func (e *ElasticsearchClient) getIDsFromIndex(querystring, typeName string, fields ...string) ([]int64, error) {
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
