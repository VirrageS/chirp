package fulltextsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = NewElasticsearch()

func TestElasticsearch_GetTweetsUsingQuerystring(t *testing.T) {
	_, err := client.GetTweetIDs("elo")
	assert.NoError(t, err)
}
