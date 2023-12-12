package api

import (
	"backend/cmd/modes/api/middlewares"
	registry "backend/cmd/registry"
	benchmark "backend/internal/repository/postgres_repo/benchmark"
	"fmt"
	"github.com/gin-gonic/gin"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

type services struct {
	Services *registry.AppServiceFields
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func SetupServer(a *registry.App) *gin.Engine {

	tp, tpErr := JaegerTraceProvider()
	if tpErr != nil {
		fmt.Println(tpErr)
		return nil
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
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

	api := router.Group("/api")
	{
		api.Use(otelgin.Middleware("backend"))
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
