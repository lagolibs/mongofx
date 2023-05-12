package examples

import (
	"git.eway.vn/x10-pushtimize/golibs/mongofx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
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
