package main

import (
	"io"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
)

func pingTime(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, time.Now().GoString())
}

func main() {
	//Create a web service
	webservice := new(restful.WebService)

	//Create route and attach it to handler
	webservice.Route(webservice.GET("/ping").To(pingTime))

	//Add service to application
	restful.Add(webservice)

	http.ListenAndServe(":8000", nil)
}
