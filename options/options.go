package options

import "go.mongodb.org/mongo-driver/mongo/options"

// ModuleConfig is a map between name of the client and client options.
type ModuleConfig map[string]*options.ClientOptions

type ModuleOptionFn func(conf ModuleConfig)

// URIs create ModuleOptionFn that parse a map of uris into ModuleConfig.
// This help
func URIs(uris map[string]string) ModuleOptionFn {
	return func(conf ModuleConfig) {
		for key, uri := range uris {
			conf[key] = options.Client().ApplyURI(uri)
		}
	}
}

func Client(name string, options *options.ClientOptions) ModuleOptionFn {
	return func(conf ModuleConfig) {
		conf[name] = options
	}
}
