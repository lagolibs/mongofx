package examples

import (
	"github.com/lagolibs/mongofx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"time"
)

func singleClient() {
	fx.New(
		mongofx.NewSimpleModule("mongo", "mongodb://localhost:27017"),
		fx.Invoke(
			fx.Annotate(func(client *mongo.Client) {},
				fx.ParamTags(`name:"mongo"`),
			),
		),
	).Run()
}

func multipleClients() {
	configs := make(map[string]string, 2)
	configs["clienta"] = "mongodb://localhost:27017/dba"
	configs["clientb"] = "mongodb://localhost:27017/dbb"

	fx.New(
		mongofx.NewModule("mongo", mongofx.WithURIs(configs)),
		fx.Invoke(
			fx.Annotate(func(client *mongo.Client) {},
				fx.ParamTags(`name:"mongo"`),
			),
		),
	).Run()
}

func timeouts() {
	configs := make(map[string]string, 2)
	configs["clienta"] = "mongodb://localhost:27017/dba"
	configs["clientb"] = "mongodb://localhost:27017/dbb"

	fx.New(
		mongofx.NewModule("mongo", mongofx.WithURIs(configs), mongofx.WithConnectTimeout(5*time.Second)),
		fx.Invoke(
			fx.Annotate(func(client *mongo.Client) {},
				fx.ParamTags(`name:"mongo"`),
			),
		),
	).Run()
}
