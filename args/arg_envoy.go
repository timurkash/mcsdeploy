package args

import (
	"fmt"
	settingsPackage "github.com/timurkash/mcsdeploy/utils/settings"
	"os"
	"text/template"
)

const (
	route = `                        - match: 
                            prefix: "/api.{{ .Name }}."
                          route:
                            cluster: {{ .NameVersion }}-cluster
                            max_grpc_timeout: 0s
`
	cluster = `    - name: {{ .NameVersion }}-cluster
      connect_timeout: 0.25s
      type: logical_dns
      http2_protocol_options: {}
      lb_policy: round_robin
      load_assignment:
        cluster_name: {{ .NameVersion }}__cluster
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: {{ .NameVersion }}
                      port_value: 9000
`
)

func ArgEnvoy() error {
	services := &settingsPackage.Services{}
	if err := services.Load(); err != nil {
		return err
	}
	fmt.Println("\n--- envoy routes")
	fmt.Println("                      routes:")
	tempRoute, err := template.New("route").Parse(route)
	if err != nil {
		return err
	}
	for _, service := range services.Services {
		if err := tempRoute.Execute(os.Stdout, service); err != nil {
			return err
		}
	}
	fmt.Println("\n --- envoy clusters")
	fmt.Println("  clusters:")
	tempCluster, err := template.New("cluster").Parse(cluster)
	if err != nil {
		return err
	}
	for _, service := range services.Services {
		if err := tempCluster.Execute(os.Stdout, service); err != nil {
			return err
		}
	}
	return nil
}
