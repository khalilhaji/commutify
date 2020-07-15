package spotify

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceive(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	ch := make(chan string)

	go receiveToken(ctx, ch)

	t.Log("resuming")
	resp, err := http.Get("http://localhost:8080/?code=test")
	assert.NoError(err)
	assert.Equal(200, resp.StatusCode)

	code, ok := <-ch
	assert.True(ok)
	assert.Equal(code, "test")
}
