package main

import (
	"github.com/pkg/errors"
	"io"
	stdlog "log"
	"net/http"

	order1 "example/internal/order"
	"example/pkg/order"

	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/julienschmidt/httprouter"
	ot "github.com/opentracing/opentracing-go"
	_ "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log"
)

func setupTracer(serviceName string) (io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://192.168.56.4:14268/api/traces",
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(log.StdLogger))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ot.SetGlobalTracer(tracer)

	return closer, nil
}

func main() {
	stdlog.SetFlags(stdlog.LstdFlags | stdlog.Lshortfile)

	closer, err := setupTracer("order")
	if err != nil {
		stdlog.Fatalf("setup tracer failed, error %+v", err)
	}
	defer func(closer io.Closer) {
		_ = closer.Close()
	}(closer)

	svc := order.NewEndpointSet(&order1.OrderSvc{})
	svc.With(opentracing.TraceEndpoint(ot.GlobalTracer(), ""))
	router := httprouter.New()
	order.Register(svc, router)
	stdlog.Println("Server now listening at https://127.0.0.1:8888/")
	// generate self-signed certificate with openssl cli
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
	err = http.ListenAndServeTLS("127.0.0.1:8888", "cert.pem", "key.pem", router)
	if err != nil {
		stdlog.Fatalf("Serve failed, error %+v", err)
	}
}
