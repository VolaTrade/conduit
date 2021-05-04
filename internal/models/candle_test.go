package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatrade/conduit/internal/models"
)

var candleResponse = []byte(`[[1620085260000,"3432.74000000","3440.00000000","3432.51000000","3438.05000000","600.03758000",1620085319999,"2062016.33871770",2046,"284.03270000","976065.00220710","0"],[1620085320000,"3438.25000000","3438.25000000","3425.00000000","3432.38000000","1709.18915000",1620085379999,"5861624.50345980",3010,"748.50941000","2566720.29482010","0"]]`)

func TestKlineUnmarshall(t *testing.T) {

	latestKline, err := models.UnmarshallBytesToLatestKline(candleResponse)
	assert.NoError(t, err)

	assert.Equal(t, latestKline.Open, 3432.74)
	assert.Equal(t, latestKline.High, 3440.0)
	assert.Equal(t, latestKline.Low, 3432.51)

}
