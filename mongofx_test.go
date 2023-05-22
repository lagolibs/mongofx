package mongofx_test

import (
	"github.com/lagolibs/mongofx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"time"
)

func ExampleNewModule() {
	configs := make(map[string]string, 2)
	configs["clienta"] = "mongodb://localhost:27017/dba"
	configs["clientb"] = "mongodb://localhost:27017/dbb"

	fx.New(
		mongofx.NewModule("mongo", mongofx.WithURIs(configs)),
		fx.Invoke(
			fx.Annotate(func(client *mongo.Client) {},
				fx.ParamTags(`name:"mongo_clienta"`),
			),
		),
	).Run()
}

func ExampleNewModule_singleClient() {
	configs := make(map[string]string, 2)
	configs["clienta"] = "mongodb://localhost:27017/dba"

	fx.New(
		mongofx.NewModule("mongo", mongofx.WithURIs(configs)),
		fx.Invoke(
			fx.Annotate(func(client *mongo.Client) {},
				fx.ParamTags(`name:"mongo_clienta"`),
			),
		),
	).Run()
}

func ExampleWithConnectTimeout() {
	configs := make(map[string]string, 2)
	configs["clienta"] = "mongodb://localhost:27017/dba"
	configs["clientb"] = "mongodb://localhost:27017/dbb"

	fx.New(
		mongofx.NewModule("mongo", mongofx.WithURIs(configs), mongofx.WithConnectTimeout(5*time.Second)),
		fx.Invoke(
			fx.Annotate(func(client *mongo.Client) {},
				fx.ParamTags(`name:"mongo_clienta"`),
			),
		),
	).Run()
}

func ExampleNewSimpleModule() {
	fx.New(
		mongofx.NewSimpleModule("mongo", "mongodb://localhost:27017"),
		fx.Invoke(
			fx.Annotate(func(client *mongo.Client) {},
				fx.ParamTags(`name:"mongo"`),
			),
		),
	).Run()
}
