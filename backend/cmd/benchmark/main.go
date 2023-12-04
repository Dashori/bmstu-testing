package main

import (
	benchmark "backend/internal/repository/postgres_repo/benchmark"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.UseRawPath = true
	router.UnescapePathValues = false

	router.Use(gin.RecoveryWithWriter(os.Stdout))
	router.Use(gin.LoggerWithWriter(os.Stdout))

	router.GET("/metrics", prometheusHandler())
	router.GET("/bench", func(ctx *gin.Context) {
		var res []string
		for i := 0; i < 2; i++ {
			fmt.Println("ITERATION ", i)
			res = benchmark.ClientBench()
		}

		ctx.JSON(http.StatusOK, res)
	})

	s := http.Server{
		Addr:    fmt.Sprintf(":8081"),
		Handler: router,
	}

	err := s.ListenAndServe()
	if err != nil && !errors.Is(http.ErrServerClosed, err) {
		panic(err)
	}

}
