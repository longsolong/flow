package ping

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePingNewGrapher(t *testing.T) {
	_, err := NewGrapher(context.Background(), []byte(`{
		"primaryRequestArgs": {
			"name": "ping",
			"version": 1
		},
		"requestArgs": {
			"hostname": "www.baidu.com",
			"timeout":  10,
			"interval": 1,
			"count":    1
		},
		"requestTags": [
			{"name": "aa", "value": "bb"}
		]
	}`))
	assert.Nil(t, err)
}
