package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	jaegerMetrics "github.com/uber/jaeger-lib/metrics"
	"log"
	"net/http"
)

func main() {
	// Configuration for testing
	cfg := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerLog.StdLogger
	jMetricsFactory := jaegerMetrics.NullFactory

	tracer, closer, err := cfg.New(
		"jaeger_test",
		jaegerConfig.Logger(jLogger),
		jaegerConfig.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Fatalf("Could not initialize tracer: %v", err.Error())
	}
	defer closer.Close()

	opentracing.InitGlobalTracer(tracer)

	r := gin.Default()
	r.GET("/", Index)
	r.Run(":20080")
}

func Index(c *gin.Context) {
	sp := opentracing.StartSpan("Index")
	defer sp.Finish()

	c.String(http.StatusOK, "🎍")
}
