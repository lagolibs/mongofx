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
				mongoClientProvider(otp),
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
func NewModule(namespace string, opts ...ModuleOptionFn) fx.Option {
	conf := moduleConfig{}
	for i := range opts {
		opts[i](conf)
	}
	return newModule(namespace, conf)
}

func newModule(namespace string, configs moduleConfig) fx.Option {
	if configs == nil || len(configs) == 0 {
		return fx.Module(namespace)
	}
	provides := make([]fx.Option, 0, len(configs))
	for name, clientOptions := range configs {
		provides = append(provides,
			fx.Provide(
				fx.Annotate(
					mongoClientProvider(clientOptions),
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

// moduleConfig is a map between name of the client and client options.
type moduleConfig map[string]*options.ClientOptions

type ModuleOptionFn func(conf moduleConfig)

// WithURIs create ModuleOptionFn that parse a map of uris into moduleConfig.
// This help integrate with configuration library such as vipers
func WithURIs(uris map[string]string) ModuleOptionFn {
	return func(conf moduleConfig) {
		for key, uri := range uris {
			conf[key] = options.Client().ApplyURI(uri)
		}
	}
}

func WithClient(name string, options *options.ClientOptions) ModuleOptionFn {
	return func(conf moduleConfig) {
		conf[name] = options
	}
}

type mongoClientConstructor func(lc fx.Lifecycle) *mongo.Client

// Actual registration logic
func mongoClientProvider(options *options.ClientOptions) mongoClientConstructor {
	return func(lc fx.Lifecycle) *mongo.Client {
		var client *mongo.Client
		lc.Append(fx.Hook{
			OnStart: func(fxCtx context.Context) error {
				ctx, cancel := context.WithTimeout(fxCtx, 10*time.Second)
				defer cancel()

				client, err := mongo.Connect(ctx, options)
				if err != nil {
					return err
				}

				ctx, cancel = context.WithTimeout(fxCtx, 5*time.Second)
				defer cancel()
				return client.Ping(ctx, readpref.Primary())
			},
			OnStop: func(ctx context.Context) error {
				return client.Disconnect(ctx)
			},
		})

		return client
	}
}
