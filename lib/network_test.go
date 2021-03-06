package tweetmap

import (
	"github.com/cdn-madness/tweetmap/protobuf"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const CONNSTR string = "ipc:///tmp/testfeed"

func TestInitPublisher(t *testing.T) {
	tweetChan := make(chan *prototweet.Tweet)
	errorChan := make(chan error)

	select {
	case err := <-errorChan:
		assert.Fail(t, err.Error())
	default:
		InitPublisher(CONNSTR, tweetChan, errorChan)
		tweet := &prototweet.Tweet{}
		tweetChan <- tweet
		close(tweetChan)
	}
}

func TestInitReceiver(t *testing.T) {
	// Start publishing tweets
	go func() {
		tweetChan := make(chan *prototweet.Tweet)
		errorChan := make(chan error)
		InitPublisher(CONNSTR, tweetChan, errorChan)
		tweet := &prototweet.Tweet{}

		for i := 0; i < 5; i++ {
			select {
			case err := <-errorChan:
				t.Logf("Error Received: %s", err)
				assert.Fail(t, "Test failed")
			case <-time.After(50 * time.Millisecond):
				tweetChan <- tweet
			}
		}
		close(tweetChan)
	}()

	// First set up the channels
	doneChan := make(chan bool)

	tweetChan, errChan := InitReceiver(CONNSTR, doneChan)
	received := false
L:
	for received == false {
		select {
		case err := <-errChan:
			assert.Fail(t, err.Error())
			break L
		case <-tweetChan:
			received = true
			close(doneChan)
		case <-time.After(1 * time.Second):
			break L
		}
	}
	assert.True(t, received)
}
