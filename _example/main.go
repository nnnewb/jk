package main

import (
	"context"
	"crypto/tls"
	"fmt"
	khttp "github.com/go-kit/kit/transport/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	stdlog "log"
	"net/http"
	"strconv"
	"strings"
	"time"

	order1 "example/internal/order"
	"example/pkg/order"

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

func OpenTelemetryTraceServer() khttp.ServerOption {
	return khttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
		tr := otel.GetTracerProvider().Tracer("OpenTelemetryTraceServer")
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(request.Header))
		ctx, span := tr.Start(ctx, request.URL.Path)
		var attrs []attribute.KeyValue
		for key, values := range request.Header {
			attrs = append(attrs, attribute.String("http.request.header."+key, strings.Join(values, "\n")))
		}
		span.SetAttributes(
			attribute.Int64("http.request_content_length", request.ContentLength),
			attribute.String("http.method", request.Method),
			attribute.String("net.protocol.name", "http"),
			attribute.String("net.protocol.version", request.Proto),
			attribute.String("net.sock.peer.addr", request.RemoteAddr),
			attribute.String("user_agent.original", request.Header.Get("User-Agent")))
		span.SetAttributes(attrs...)
		return ctx
	})
}

func OpenTelemetryTraceServerResp() khttp.ServerOption {
	return khttp.ServerAfter(func(ctx context.Context, wr http.ResponseWriter) context.Context {
		span := trace.SpanFromContext(ctx)
		var attrs []attribute.KeyValue
		for key, values := range wr.Header() {
			attrs = append(attrs, attribute.String("http.response.header."+key, strings.Join(values, "\n")))
		}
		span.SetAttributes(attrs...)
		return ctx
	})
}

func OpenTelemetryTraceServerEnd() khttp.ServerOption {
	return khttp.ServerFinalizer(func(ctx context.Context, code int, req *http.Request) {
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.Int("http.status_code", code))
		if span != nil {
			span.End()
		}
	})
}

func OpenTelemetryTraceClient() khttp.ClientOption {
	return khttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
		tr := otel.GetTracerProvider().Tracer("OpenTelemetryTraceClient")
		ctx, span := tr.Start(ctx, request.Method+" "+request.URL.String())
		var attrs []attribute.KeyValue
		for key, values := range request.Header {
			attrs = append(attrs, attribute.String("http.request.header."+key, strings.Join(values, "\n")))
		}
		span.SetAttributes(attrs...)
		var port int
		portStr := request.URL.Port()
		if portStr != "" {
			port, _ = strconv.Atoi(portStr)
		}
		span.SetAttributes(
			attribute.String("http.method", request.Method),
			attribute.String("http.flavor", fmt.Sprintf("%d.%d", request.ProtoMajor, request.ProtoMinor)),
			attribute.String("http.url", request.URL.String()),
			attribute.String("net.sock.peer.name", request.URL.Hostname()),
			attribute.Int("net.sock.peer.port", port))

		propagator := otel.GetTextMapPropagator()
		propagator.Inject(ctx, propagation.HeaderCarrier(request.Header))

		return ctx
	})
}

func OpenTelemetryTraceClientResp() khttp.ClientOption {
	return khttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.Int("http.status_code", response.StatusCode))
		return ctx
	})
}

func OpenTelemetryTraceClientEnd() khttp.ClientOption {
	return khttp.ClientFinalizer(func(ctx context.Context, err error) {
		span := trace.SpanFromContext(ctx)
		span.RecordError(err)
		span.End()
	})
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
		OpenTelemetryTraceServer(),
		OpenTelemetryTraceServerResp(),
		OpenTelemetryTraceServerEnd())

	clientSet := order.NewHTTPClientSet(
		"https",
		"127.0.0.1",
		8888,
		OpenTelemetryTraceClient(),
		OpenTelemetryTraceClientResp(),
		OpenTelemetryTraceClientEnd(),
		khttp.SetClient(&http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}))
	go func() {
		time.Sleep(5 * time.Second)
		_, err := clientSet.EndpointSet().CreateOrder(context.Background(), order.CreateOrderRequest{})
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
