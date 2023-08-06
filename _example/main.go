package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"example/api/order"
	order1 "example/internal/order"

	"github.com/gin-gonic/gin"
	klog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	okmgin "github.com/nnnewb/otelkit/metric/gin"
	oktgin "github.com/nnnewb/otelkit/tracing/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	_ "github.com/uber/jaeger-client-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName("example"))),
	)
	return tp, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tp, err := tracerProvider("http://192.168.56.4:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	promExporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(promExporter))
	meter := provider.Meter("http-example")

	// serving /metrics endpoint
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("prometheus scrap endpoint start serving at https://127.0.0.1:23333/metrics")
		err := http.ListenAndServeTLS("127.0.0.1:23333", "secrets/cert.pem", "secrets/key.pem", http.DefaultServeMux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// HOWTO: generate self-signed certificate with openssl cli
	//
	// openssl req \
	//     -x509 \
	//     -newkey rsa:4096 \
	//     -keyout key.pem \
	//     -out cert.pem \
	//     -sha256 \
	//     -days 3650 \
	//     -nodes \
	//     -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname"

	errorLogger := klog.NewLogfmtLogger(os.Stdout)
	errorLogger = klog.With(errorLogger, "caller", klog.DefaultCaller, "timestamp", klog.DefaultTimestamp)
	errorLogger = level.NewFilter(errorLogger, level.AllowDebug())
	errorLogger = level.NewInjector(errorLogger, level.ErrorValue())

	endpointSet := order.NewEndpointSet(&order1.OrderSvc{})
	serverSet := order.NewGinServerSet(endpointSet)
	engine := gin.New()
	engine.Use(oktgin.TraceMiddleware(tp.Tracer("oktgin"), otel.GetTextMapPropagator()))
	engine.Use(okmgin.MeasureHandleFunc(meter))
	serverSet.Register(engine)
	serverSet.RegisterEmbedSwaggerUI(engine)

	log.Println("Server now listening at https://127.0.0.1:8888/")
	err = http.ListenAndServeTLS("127.0.0.1:8888", "secrets/cert.pem", "secrets/key.pem", engine)
	if err != nil {
		log.Fatalf("Serve failed, error %+v", err)
	}
}
