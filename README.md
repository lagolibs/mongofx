# Mongo Fx

Fx Module for mongo client

### Usage

The recommended way to config mongo client is using uri

```
mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
```

Example usage with multiple mongodb and viper

```properties
mongodb.uris.clienta = mongodb://localhost:27017/dba
mongodb.uris.clientb = mongodb://localhost:27017/dbb
```

```go
package main

import (
	"git.eway.vn/x10-pushtimize/golibs/mongofx"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"os"
)

func init() {
	viper.AddConfigPath(lo.Must(os.Getwd()))
	viper.SetConfigType("properties")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.SafeWriteConfig(); err != nil {
		lo.Must0(viper.ReadInConfig())
	}
}

func main() {
	app := fx.New(
		mongofx.NewModule("mongo", mongofx.WithURIs(viper.GetStringMapString("mongodb"))),
		fx.Invoke(fx.Annotate(func(client *mongo.Client, client2 *mongo.Client) {}, fx.ParamTags(`name:"mongo_clienta"`, `name:"mongo_clienb"`))),
	)

	app.Run()
}

```

See [examples/](examples/examples.go) for more usage.

