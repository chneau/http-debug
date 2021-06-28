package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Debug struct {
	Hostname       string
	Time           time.Time
	ClientIP       string
	RemoteIP       string
	RequestHeaders map[string]string
}

func NewDebug(c *gin.Context, hostname string) Debug {
	remoteIP, _ := c.RemoteIP()
	return Debug{
		ClientIP:       c.ClientIP(),
		RemoteIP:       remoteIP.String(),
		Time:           time.Now(),
		Hostname:       hostname,
		RequestHeaders: flattenHeaders(c.Request.Header),
	}
}

func flattenHeaders(headers map[string][]string) map[string]string {
	result := map[string]string{}
	for k, v := range headers {
		valueStr := strings.Join(v, ",")
		result[k] = valueStr
	}
	return result
}

func main() {
	router := gin.Default()

	environmentVariables := buildEnvironmentVariables()

	hostname, _ := os.Hostname()

	router.GET("/favicon.ico", func(_ *gin.Context) {})
	router.GET("/env", func(c *gin.Context) { c.IndentedJSON(http.StatusOK, environmentVariables) })
	router.NoRoute(func(c *gin.Context) { c.IndentedJSON(http.StatusOK, NewDebug(c, hostname)) })

	err := router.Run()
	if err != nil {
		log.Panicln(err)
	}
}

func buildEnvironmentVariables() map[string]string {
	result := map[string]string{}
	for _, variable := range os.Environ() {
		splits := strings.SplitN(variable, "=", 2)
		result[splits[0]] = splits[1]
	}
	return result
}
