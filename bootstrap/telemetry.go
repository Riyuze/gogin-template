package bootstrap

import (
	"os"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func (c *Container) GetTracer() trace.Tracer {
	return otel.Tracer(os.Getenv("SERVICE_NAME"))
}

func (c *Container) initTracer() *sdktrace.TracerProvider {
	if !c.GetConfig().GetBool("telemetry.enable") {
		c.logrus.Debug("opentelemetry disabled")
		return nil
	}

	c.logrus.Debug("opentelemetry initialize")

	host := c.GetConfig().GetString("telemetry.jaeger.agent_host")
	port := c.GetConfig().GetString("telemetry.jaeger.agent_port")

	c.logrus.
		WithField("bootstrap", "jaeger").
		Debugf("trying to connect to %s:%s", host, port)

	exporter, err := otlptracehttp.New(c.ctx, otlptracehttp.WithEndpoint((host + ":" + port)))
	if err != nil {
		c.logrus.Fatal(err)
	}

	c.logrus.
		WithField("bootstrap", "jaeger").
		Debugf("connected to %s:%s", host, port)

	c.trace = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(os.Getenv("SERVICE_NAME")),
			attribute.String("environment", os.Getenv("ENV")),
		)),
	)

	otel.SetTracerProvider(c.trace)
	otel.SetTextMapPropagator(
		b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader | b3.B3SingleHeader)),
	)

	return c.trace
}
