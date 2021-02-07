package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func checkRegion(c *gin.Context) (string, error) {
	// TODO: Add logic here
	return euRegion, nil
}

func HandleAnyRoute(c *gin.Context) {
	path := c.Param("path")
	clientMethod := c.Request.Method

	getRegionUrl, err := checkRegion(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Region not found!"})
		return
	}

	internalUrl := fmt.Sprint("http://", getRegionUrl, "/", path)
	switch clientMethod {
	case get:
		// Let us make request on behalf of Client to the valid endpoint
		resp, err := http.Get(internalUrl)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		// Then let us send it to the Client
		c.JSON(200, body)
	case post:
		// TODO: Implement me
		fmt.Println("In POST method")
		c.JSON(200, gin.H{"message": "Not implemented"})
	case put:
		// TODO: Implement me
		fmt.Println("In PUT method")
		c.JSON(200, gin.H{"message": "Not implemented"})
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "Method not allowed"})
	}
}
