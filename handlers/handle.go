package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func HandleAnyRoute(c *gin.Context) {
	path := c.Param("path")
	remoteUrl := "https://pudelek.pl"
	parsedUrl := fmt.Sprint(remoteUrl,"/", path)
	remote, err := url.Parse(parsedUrl)
	if err != nil {
		panic(err)
	}
	clientMethod := c.Request.Method

	log.Print("Preparing PROXY...")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
	}
	log.Printf("Configuring underlying request for HTTP %v method", clientMethod)
	log.Printf("Making underlying request...")
	proxy.ServeHTTP(c.Writer, c.Request)
	log.Printf("Request was ended without issues")
	log.Printf("Underying HTTP staus code: %v", c.Writer.Status())
}
