package httpengine

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"encoding/json"
)

func NewHTTPEngine(apiVersion string) *HTTPEngine {
	httpEngine := HTTPEngine{APIVersion: apiVersion}
	return &httpEngine
}

type HTTPEngine struct {
	APIVersion string
	Server     *http.Server
	Router     *httprouter.Router
}

func (httpEngine *HTTPEngine) PowerUp(host string, port int) {
	httpEngine.Router = httprouter.New()
	httpEngine.Router.GET("/api/version", httpEngine.apiVersionCheckHandler)

	httpEngine.Server = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	fmt.Printf("Http server listen on %v, port:%v \n", host, port)

	httpEngine.Server.Handler = httpEngine.Router
	httpEngine.Server.ListenAndServe()
}

func (httpEngine *HTTPEngine) apiVersionCheckHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	data := map[string]string{"apiVersion": httpEngine.APIVersion}
	encodedData, _ := json.Marshal(data)

	response.Header().Set("content-type", "application/javascript")
	_, err := response.Write(encodedData)
	if err != nil {
		fmt.Print(err)
	}
}
