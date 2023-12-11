package api

import (
	"backend/cmd/modes/api/middlewares"
	registry "backend/cmd/registry"
	benchmark "backend/internal/repository/postgres_repo/benchmark"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type services struct {
	Services *registry.AppServiceFields
}

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	responseLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_latency_seconds",
			Help:    "Response latency of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	inFlightRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_in_flight_requests",
		Help: "Number of in-flight HTTP requests",
	})
)

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func SetupServer(a *registry.App) *gin.Engine {
	t := services{a.Services}

	router := gin.Default()
	router.GET("/metrics", prometheusHandler())
	router.GET("/bench", func(ctx *gin.Context) {
		var res [][]string
		for i := 0; i < 10; i++ {
			fmt.Println("ITERATION ", i)
			res2 := benchmark.ClientBench()
			res = append(res, res2)
		}

		ctx.JSON(http.StatusOK, res)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inFlightRequests.Inc()

		status := http.StatusOK
		defer func() {
			duration := time.Since(start).Seconds()
			httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, string(status)).Inc()
			responseLatency.WithLabelValues(r.Method, r.URL.Path, string(status)).Observe(duration)
			inFlightRequests.Dec()
		}()

		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(status)
	})

	api := router.Group("/api")
	{
		api.POST("/setRole", t.setRole)

		api.GET("/doctors", t.getAllDoctors)
		api.POST("/doctor/create", t.createDoctor)
		api.POST("/doctor/login", t.loginDoctor)

		doctor := api.Group("/doctor")
		doctor.Use(middlewares.JwtAuthMiddleware())
		doctor.GET("/info", t.doctorInfo)
		doctor.GET("/records", t.doctorRecords)
		doctor.PATCH("/shedule", t.doctorShedule)

		api.POST("/client/create", t.createClient)
		api.POST("/client/createOTP", t.createClientOTP)
		api.POST("/client/login", t.loginClient)

		client := api.Group("/client")
		client.Use(middlewares.JwtAuthMiddleware())
		client.GET("/info", t.infoClient)
		client.GET("/records", t.ClientRecords)
		client.GET("/pets", t.ClientPets)
		client.POST("/record", t.NewRecord)
		client.POST("/pet", t.NewPet)
		client.DELETE("/pet", t.DeletePet)
	}

	port := a.Config.Port
	adress := a.Config.Address
	err := router.Run(adress + port)

	if err != nil {
		return nil
	}

	return router
}

type Role struct {
	Role string
}

func (t *services) setRole(c *gin.Context) {
	var role *Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		jsonInternalServerErrorResponse(c, err)
		return
	}

	if role.Role == "doctor" {
		err = t.Services.DoctorService.SetRole()
	} else if role.Role == "client" {
		err = t.Services.ClientService.SetRole()
	} else {
		jsonBadRequestResponse(c, fmt.Errorf("Такой роли не существует!"))
	}

	if err != nil {
		jsonInternalServerErrorResponse(c, err)
		return
	}

	jsonStatusOkResponse(c)

}
