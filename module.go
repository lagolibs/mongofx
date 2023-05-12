package mongofx

import (
	"context"
	"fmt"
	"git.eway.vn/x10-pushtimize/golibs/mongofx/options"
	"go.mongodb.org/mongo-driver/mongo"
	motp "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
	"time"
)

// NewSimpleModule construct a module contain single client.
// Does not register group namespace.
// The name of the mongo client is the same as the name space.
func NewSimpleModule(namespace string, uri string) fx.Option {
	otp := motp.Client().ApplyURI(uri)
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
func NewModule(namespace string, opts ...options.ModuleOptionFn) fx.Option {
	conf := options.ModuleConfig{}
	for i := range opts {
		opts[i](conf)
	}
	return newModule(namespace, conf)
}

func newModule(namespace string, configs options.ModuleConfig) fx.Option {
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

type mongoClientConstructor func(lc fx.Lifecycle) *mongo.Client

func mongoClientProvider(options *motp.ClientOptions) mongoClientConstructor {
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
