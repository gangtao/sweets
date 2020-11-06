package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gangtao/sweets/config"
)

var kvClient config.KVStore

func GetConfig(c *gin.Context) {
	//tenant := c.Query("tenant") TODO: support tenant in the future
	dataId := c.Query("dataId")
	group := c.Query("group")

	val, err := kvClient.GetConfig(dataId, group)
	if err != nil {
		panic(err)
	}

	c.String(http.StatusOK, val)
}

func PublishConfig(c *gin.Context) {
	//tenant := c.Query("tenant") TODO: support tenant in the future
	dataId := c.Query("dataId")
	group := c.Query("group")
	content := c.Query("content")

	err := kvClient.PublishConfig(dataId, group, content)
	if err != nil {
		c.String(http.StatusOK, "false")
	} else {
		c.String(http.StatusOK, "true")
	}
}

func DeleteConfig(c *gin.Context) {
	dataId := c.Query("dataId")
	group := c.Query("group")

	err := kvClient.DeleteConfig(dataId, group)
	if err != nil {
		c.String(http.StatusOK, "false")
	} else {
		c.String(http.StatusOK, "true")
	}
}

func main() {
	kvClient = config.NewEtcdStore()

	router := gin.Default()

	router.GET("/sweets/v1/cs/configs", GetConfig)
	router.POST("/sweets/v1/cs/configs", PublishConfig)
	router.DELETE("/sweets/v1/cs/configs", DeleteConfig)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}