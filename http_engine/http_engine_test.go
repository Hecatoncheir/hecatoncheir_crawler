package http_engine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var once sync.Once

func setUpHTTPServer() {
	httpServer := NewHTTPEngine("v1.0")
	httpServer.PowerUp("localhost", 8181)
}

func TestHttpEngineCanSendVersionOfAPI(test *testing.T) {
	var err error
	go once.Do(setUpHTTPServer)

	respose, err := http.Get("http://localhost:8181/api/version")
	if err != nil {
		test.Log(err)
		test.Fail()
	}

	encodedBody, err := ioutil.ReadAll(respose.Body)
	if err != nil {
		test.Log(err)
		test.Fail()
	}

	decodedBody := map[string]string{}

	err = json.Unmarshal(encodedBody, &decodedBody)
	if err != nil {
		test.Log(err)
		test.Fail()
	}

	assert.Equal(test, decodedBody["apiVersion"], "v1.0", "The api versions should be the same.")
}
