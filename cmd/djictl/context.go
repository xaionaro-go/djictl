package main

import (
	"context"
	"os"
	"os/user"
	"strings"

	"github.com/facebookincubator/go-belt"
	xruntime "github.com/facebookincubator/go-belt/pkg/runtime"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics"
	prometheusadapter "github.com/facebookincubator/go-belt/tool/experimental/metrics/implementation/prometheus"
	"github.com/facebookincubator/go-belt/tool/logger"
	xlogrus "github.com/facebookincubator/go-belt/tool/logger/implementation/logrus"
	"github.com/sirupsen/logrus"
	"github.com/xaionaro-go/observability"
	"github.com/xaionaro-go/secret"
)

const (
	AppName = "djictl"
)

var originalPCFilter xruntime.PCFilter

func init() {
	originalPCFilter = xruntime.DefaultCallerPCFilter
}

func setDefaultCallerPCFilter() {
	xruntime.DefaultCallerPCFilter = observability.CallerPCFilter(originalPCFilter)
}

func getContext(
	loggerLevel logger.Level,
	insecureDebug bool,
	logstashAddr string,
) context.Context {
	ctx := context.Background()
	setDefaultCallerPCFilter()

	ctx = metrics.CtxWithMetrics(ctx, prometheusadapter.Default())

	secretsProvider := observability.NewStaticSecretsProvider()
	ctx = observability.WithSecretsProvider(ctx, secretsProvider)

	ll := xlogrus.DefaultLogrusLogger()
	ll.Formatter.(*logrus.TextFormatter).ForceColors = true

	logPreHooks := logger.PreHooks{}
	logPreHooks = append(logPreHooks,
		observability.StructFieldSecretsFilter{},
		observability.NewSecretValuesFilter(secretsProvider),
	)

	logHooks := logger.Hooks{}
	if insecureDebug {
		secret.SetSecrecy(false)
	} else {
		logHooks = append(logHooks,
			observability.NewRemoveInsecureDebugFilter(),
		)
	}

	l := xlogrus.New(ll).WithLevel(loggerLevel).WithPreHooks(logPreHooks...).WithHooks(logHooks...)
	ctx = logger.CtxWithLogger(ctx, l)

	if logstashAddr != "" {
		ctx = observability.CtxWithLogstash(
			ctx,
			logstashAddr,
			strings.ToLower(AppName),
		)
	}

	ctx = belt.WithField(ctx, "program", strings.ToLower(AppName))

	if hostname, err := os.Hostname(); err == nil {
		ctx = belt.WithField(ctx, "hostname", strings.ToLower(hostname))
	}

	ctx = belt.WithField(ctx, "uid", os.Getuid())
	ctx = belt.WithField(ctx, "pid", os.Getpid())

	if u, err := user.Current(); err == nil {
		ctx = belt.WithField(ctx, "user", u.Username)
	}

	l = logger.FromCtx(ctx)
	logger.Default = func() logger.Logger {
		return l
	}

	return ctx
}
