package app

import (
	"io"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// Tracer is a facade
type Tracer struct {
	Ctx opentracing.Tracer
	io  io.Closer
}

func (t Tracer) Close() {
	t.io.Close()
}

func NewTracer() *Tracer {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: "localhost:5775",
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	var t Tracer
	tracer, closer, err := cfg.New(
		"openidservice", jaegercfg.Logger(jaeger.StdLogger))

	t.Ctx = tracer
	t.io = closer

	if err != nil {
		log.Fatal(err)
	}
	return &t
}

func gotrace(t opentracing.Tracer) {

	span := t.StartSpan("new_span")
	defer span.Finish()
	span.SetOperationName("span_1")
	span.LogFields(tracelog.String("ds", "asd"))
	span.LogEvent("hello")
	span.SetTag(zipkincore.HTTP_PATH, struct{ name string }{"ad"})

	span.SetBaggageItem("Some_Key", "12345")
	span.SetBaggageItem("Some-other-key", "42")
}
