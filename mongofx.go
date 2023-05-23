package mongofx

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
	"time"
)

// NewSimpleModule construct a module contain single client.
// Does not register group namespace.
// The name of the mongo client is the same as the name space.
func NewSimpleModule(namespace string, uri string) fx.Option {
	otp := options.Client().ApplyURI(uri)
	return fx.Module(namespace,
		fx.Provide(
			fx.Annotate(
				mongoClientProvider(otp, newDefaultTimeout()),
				fx.ResultTags(
					fmt.Sprintf(`name:"%s"`, namespace),
				),
			),
		),
	)
}

// NewModule construct a new fx Module for mongodb, using configuration options
// Each mongo client will be named as <namespace>_<name>
// Also register a <namespace> group
func NewModule(namespace string, opts ...ModuleOption) fx.Option {
	conf := moduleConfig{
		timeout: newDefaultTimeout(),
		configs: make(map[string]*options.ClientOptions, len(opts)),
	}
	for i := range opts {
		opts[i](&conf)
	}
	return newModule(namespace, conf)
}

func newModule(namespace string, conf moduleConfig) fx.Option {
	configs := conf.configs
	if configs == nil || len(configs) == 0 {
		return fx.Module(namespace)
	}
	provides := make([]fx.Option, 0, len(configs))
	for name, clientOptions := range configs {
		provides = append(provides,
			fx.Provide(
				fx.Annotate(
					mongoClientProvider(clientOptions, conf.timeout),
					fx.ResultTags(
						fmt.Sprintf(`name:"%s_%s"`, namespace, name),
						fmt.Sprintf(`group:"%s"`, namespace),
					),
				),
			),
		)
	}
	return fx.Module(namespace, provides...)
}

type timeoutConfig struct {
	connectTimeout time.Duration
	pingTimeout    time.Duration
}

func newDefaultTimeout() timeoutConfig {
	return timeoutConfig{
		pingTimeout:    10 * time.Second,
		connectTimeout: 10 * time.Second,
	}
}

type moduleConfig struct {
	configs map[string]*options.ClientOptions
	timeout timeoutConfig
}

// ModuleOption applies an option to moduleConfig
type ModuleOption func(conf *moduleConfig)

// WithURIs create ModuleOption that parse a map of uris into moduleConfig.
// This help integrate with configuration library such as vipers
func WithURIs(uris map[string]string) ModuleOption {
	return func(conf *moduleConfig) {
		for key, uri := range uris {
			conf.configs[key] = options.Client().ApplyURI(uri)
		}
	}
}

// WithConnectTimeout set the timeout for client initialization.
// The timeout for Connect and Ping operations will be half of given totalTimeout for each.
func WithConnectTimeout(totalTimeout time.Duration) ModuleOption {
	opTimeout := totalTimeout / 2
	return func(conf *moduleConfig) {
		conf.timeout.connectTimeout = opTimeout
		conf.timeout.pingTimeout = opTimeout
	}
}

func WithClient(name string, options *options.ClientOptions) ModuleOption {
	return func(conf *moduleConfig) {
		conf.configs[name] = options
	}
}

type mongoClientConstructor func(lc fx.Lifecycle) (*mongo.Client, error)

// Actual registration logic
func mongoClientProvider(options *options.ClientOptions, config timeoutConfig) mongoClientConstructor {
	return func(lc fx.Lifecycle) (*mongo.Client, error) {
		ctx, cancel := context.WithTimeout(context.Background(), config.connectTimeout)
		defer cancel()

		client, err := mongo.Connect(ctx, options)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStart: func(fxCtx context.Context) error {
				ctx, cancel = context.WithTimeout(fxCtx, config.pingTimeout)
				defer cancel()
				return client.Ping(ctx, readpref.Primary())
			},
			OnStop: func(ctx context.Context) error {
				return client.Disconnect(ctx)
			},
		})

		return client, nil
	}
}
