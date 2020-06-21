package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shukla2112/counter-app/config"
	"github.com/shukla2112/counter-app/handler"
)

// Route :
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

// Routes :
type Routes []Route

// NewRouter :
func NewRouter(useSyslog bool, appC *config.AppConfig) *gin.Engine {
	var routes = Routes{
		Route{"GET", "/", IndexHandler},

		// Counter App : endpoints
		Route{"GET", "/counter/:key", handler.GetCounter(appC)},
		Route{"POST", "/counter/:key", handler.SetCounter(appC)},
		// Route{"POST", "/counter/:key", handler.EndpointHandler(operator.NewOpCounterInput, operator.OpCounter, appC)},

	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(noColorLogger())
	for _, route := range routes {
		router.Handle(route.Method, route.Pattern, route.HandlerFunc)
	}
	return router
}

// StartServer : Start web server to receive sites-db data updates from sources
func StartServer(ctx context.Context, appC *config.AppConfig) {
	router := NewRouter(false, appC)
	s := &http.Server{Addr: appC.ConfigData.Server, Handler: router}
	go func() { log.Fatalln(s.ListenAndServe()) }()

	<-ctx.Done()
	log.Println("HTTP_SERVER: stopping server...")
	s.Shutdown(nil)
}

// IndexHandler :
func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the counter-app!",
	})
}

//this logger implements the same request handler as by default, but
//removes colors from the output making it suitable for logging to
//syslog
func noColorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		// Process request
		c.Next()
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		fmt.Fprintf(gin.DefaultWriter,
			"[GIN] %v | %3d | %13v | %s | %-7s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			comment,
		)
	}
}
