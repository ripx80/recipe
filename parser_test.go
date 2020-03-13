package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestM3Api(t *testing.T) {
	_, err := LoadFile("testdata/apiTest.json", &M3{})
	assert.Nil(t, err)
}

func TestM3FermentationHop(t *testing.T) {
	_, err := LoadFile("testdata/fermentationHop.json", &M3{})
	assert.Nil(t, err)
}

func TestM3Whirlpool(t *testing.T) {
	_, err := LoadFile("testdata/whirlpool.json", &M3{})
	assert.Nil(t, err)
}

func TestM3HopsHoney(t *testing.T) {
	_, err := LoadFile("testdata/hopsHoney.json", &M3{})
	assert.Nil(t, err)
}
