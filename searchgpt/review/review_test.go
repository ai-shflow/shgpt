//go:build review_test

// go test -cover -covermode=atomic -parallel 2 -tags=review_test -v searchgpt/review

package review

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"searchgpt/config"
)

const (
	changeGerrit = 42
	commitGerrit = "aed052e8b66810795d3a894a9095db41e5854b70"
	queryAfter   = "after:2024-08-01"
	queryBefore  = "before:2024-11-01"
)

func initReview(t *testing.T) review {
	var cfg *Config

	file, _ := config.ConfigFile.ReadFile(configFile)
	_ = yaml.Unmarshal(file, &cfg)

	return review{cfg: cfg}
}

func TestAccount(t *testing.T) {
	// TBD: FIXME
	assert.Equal(t, nil, nil)
}

func TestQuery(t *testing.T) {
	var buf []interface{}
	var err error
	var ret []byte

	r := initReview(t)

	buf, err = r.Query("change:"+strconv.Itoa(changeGerrit), 0, 10)
	assert.Equal(t, nil, err)

	ret, _ = json.Marshal(buf)
	fmt.Printf("change: %s\n", string(ret))

	buf, err = r.Query(queryAfter+" "+queryBefore, 0, 10)
	assert.Equal(t, nil, err)

	ret, _ = json.Marshal(buf)
	fmt.Printf("change: %s\n", string(ret))
}

func TestGetAccount(t *testing.T) {
	// TBD: FIXME
	assert.Equal(t, nil, nil)
}

func TestGetQuery(t *testing.T) {
	h := initReview(t)

	_, err := h.get(h.urlQuery("commit:-1", []string{"CURRENT_REVISION"}, 0))
	assert.NotEqual(t, nil, err)

	buf, err := h.get(h.urlQuery("commit:"+commitGerrit, []string{"CURRENT_REVISION"}, 0))
	assert.Equal(t, nil, err)

	_, err = h.unmarshalList(buf)
	assert.Equal(t, nil, err)
}
