package main

import (
	"context"
	"crypto/tls"
	stdlog "log"
	"net/http"
	"time"

	khttp "github.com/go-kit/kit/transport/http"
	"github.com/nnnewb/otelkit"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"

	"example/api/order"
	order1 "example/internal/order"

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
	stdlog.SetFlags(stdlog.LstdFlags | stdlog.Lshortfile)
	tp, err := tracerProvider("http://192.168.56.4:14268/api/traces")
	if err != nil {
		stdlog.Fatal(err)
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
			stdlog.Fatal(err)
		}
	}(ctx)

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

	endpointSet := order.NewEndpointSet(&order1.OrderSvc{})
	serverSet := order.NewHTTPServerSet(
		endpointSet,
		otelkit.OpenTelemetryTraceServer(),
		otelkit.OpenTelemetryTraceServerResp(),
		otelkit.OpenTelemetryTraceServerEnd())

	clientSet := order.NewHTTPClientSet(
		"https",
		"127.0.0.1",
		8888,
		otelkit.OpenTelemetryTraceClient(),
		otelkit.OpenTelemetryTraceClientResp(),
		otelkit.OpenTelemetryTraceClientEnd(),
		khttp.SetClient(&http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}))
	go func() {
		time.Sleep(5 * time.Second)
		_, err := clientSet.EndpointSet().CreateOrder(context.Background(), &order.CreateOrderRequest{})
		if err != nil {
			stdlog.Printf("create order failed, error %+v", err)
		}
	}()

	stdlog.Println("Server now listening at https://127.0.0.1:8888/")
	err = http.ListenAndServeTLS("127.0.0.1:8888", "cert.pem", "key.pem", serverSet.Handler())
	if err != nil {
		stdlog.Fatalf("Serve failed, error %+v", err)
	}
}
