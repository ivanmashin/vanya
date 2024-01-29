//go:build vanya
// +build vanya

package single_obj

import "github.com/ivanmashin/vanya"

func main() {
	vanya.BuildConfigs(
		MyConfig{
			HttpEndpoint:       "localhost:51000",
			GrpcEndpoint:       "localhost:52000",
			MonitoringEndpoint: "localhost:53000",
		},
	)
}

type MyConfig struct {
	HttpEndpoint       string
	GrpcEndpoint       string
	MonitoringEndpoint string
}
