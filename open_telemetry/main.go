package main

//
//import (
//	"context"
//	"crypto/tls"
//	"fmt"
//	"go.opentelemetry.io/otel"
//	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
//	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
//	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
//	"go.opentelemetry.io/otel/propagation"
//	"go.opentelemetry.io/otel/sdk/resource"
//	sdktrace "go.opentelemetry.io/otel/sdk/trace"
//	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials"
//	"google.golang.org/grpc/credentials/insecure"
//	"os"
//	"time"
//)
//
//type TraceParam struct {
//	ServiceName     string  // 服务名称
//	ServiceVersion  string  // 版本号，避免服务版本不一致问题
//	ServiceInstance string  // 示例标识，pod id或者ip:port之类的
//	Environment     string  // dev, test, prod之类
//	Endpoint        string  // host:port
//	Authorization   string  // basic auth或者api key
//	Protocol        string  // http或者grpc
//	EnableTLS       bool    //是否使用ssl
//	CertFile        string  //证书路径，使用tls时才需要配置
//	SampleRate      float64 //默认为1
//}
//
//// InitProvider 初始化并返回provider
//func InitProvider(param TraceParam, asGlobal bool) (*sdktrace.TracerProvider, error) {
//	if param.ServiceName == "" || param.ServiceVersion == "" {
//		return nil, fmt.Errorf("invalid service param to init tracer")
//	}
//	//默认是生产环境
//	if param.Environment == "" {
//		param.Environment = "prod"
//	}
//	ctx := context.Background()
//	res, err := resource.New(ctx,
//		resource.WithAttributes(
//			semconv.ServiceNameKey.String(param.ServiceName),
//			semconv.ServiceVersionKey.String(param.ServiceVersion),
//			semconv.ServiceInstanceIDKey.String(param.ServiceInstance),
//			semconv.DeploymentEnvironmentKey.String(param.Environment),
//		),
//		resource.WithHost(),
//		resource.WithProcess(),
//		resource.WithTelemetrySDK(),
//		resource.WithSchemaURL(semconv.SchemaURL),
//	)
//	if err != nil {
//		return nil, fmt.Errorf("fail to create resource:%w", err)
//	}
//	var exporter sdktrace.SpanExporter
//	if param.Endpoint == "" {
//		//兜底策略，控制台输出
//		exporter, err = stdouttrace.New(
//			stdouttrace.WithWriter(os.Stdout),
//			stdouttrace.WithPrettyPrint(),
//		)
//	} else if param.Protocol == ProtocolGRPC {
//		//连接collector超时时间
//		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
//		defer cancel()
//		opts := []grpc.DialOption{grpc.WithBlock()}
//		if !param.EnableTLS {
//			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
//		} else {
//			//tls证书
//			var creds credentials.TransportCredentials
//			if param.CertFile != "" {
//				creds, err = credentials.NewClientTLSFromFile(param.CertFile, "")
//				if err != nil {
//					return nil, fmt.Errorf("init grpc conn to apm server err:%s", err)
//				}
//			} else {
//				return nil, fmt.Errorf("you should specific CertFile when enable tls with grpc")
//			}
//			opts = append(opts, grpc.WithTransportCredentials(creds))
//		}
//		conn, err := grpc.DialContext(ctx, param.Endpoint, opts...)
//		if err != nil {
//			log.Fatal("fail to init grpc conn to apm server:%s", err)
//		}
//		exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
//		if err != nil {
//			return nil, fmt.Errorf("fail to create trace exporter:%w", err)
//		}
//	} else {
//		var tlsConf otlptracehttp.Option
//		if param.EnableTLS {
//			tlsConf = otlptracehttp.WithTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
//		} else {
//			tlsConf = otlptracehttp.WithInsecure()
//		}
//		exporter, err = otlptracehttp.New(context.Background(),
//			otlptracehttp.WithEndpoint(param.Endpoint),
//			tlsConf,
//		)
//		if err != nil {
//			return nil, fmt.Errorf("fail to create http exporter:%w", err)
//		}
//	}
//	bsp := sdktrace.NewBatchSpanProcessor(exporter)
//	traceProvider := sdktrace.NewTracerProvider(
//		sdktrace.WithSampler(
//			sdktrace.ParentBased(sdktrace.TraceIDRatioBased(param.SampleRate)),
//		),
//		sdktrace.WithResource(res),
//		sdktrace.WithSpanProcessor(bsp),
//	)
//	if asGlobal {
//		otel.SetTracerProvider(traceProvider)
//	}
//	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
//		propagation.TraceContext{}, propagation.Baggage{}))
//	return traceProvider, nil
//}
