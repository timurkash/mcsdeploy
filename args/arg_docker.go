package args

import (
	"fmt"
	settingsPackage "github.com/timurkash/mcsdeploy/utils/settings"
	"os"
	"text/template"
)

const (
	service = `
  {{ .NameVersion }}:
    image: registry.{{ .ProjectRepo }}/back/{{ .NameVersion }}:v0.0.1
    build:
      context: ../../back/{{ .NameVersion }}
      dockerfile: Dockerfile
    restart: always
    env_file:
      - ./.env
    volumes:
      - ./configs/{{ .NameVersion }}:/data/conf
#    ports:
#      - "{{ .Port }}:9000"
#      - "{{ .HttpPort }}:8000"
    depends_on:
#      - postgres
      - jaeger
    networks:
      - backend
`
	//	config = `
	//server:
	//  grpc:
	//    addr: :9000
	//    timeout: 1s
	//  http:
	//    addr: :8000
	//    timeout: 1s
	//jwks:
	//  url: ${JWKS_URL}
	//  refresh_interval: 3600s
	//  refresh_rate_limit: 300s
	//  refresh_timeout: 10s
	//business:
	//data:
	//  relational:
	//    dialect: postgres
	//    host: ${POSTGRES_HOST}
	//    port: ${POSTGRES_PORT}
	//    user: ${POSTGRES_USER}
	//    password: ${POSTGRES_PASSWORD}
	//    dbname: ${POSTGRES_DB}
	//    schema: {{ .Name }}
	//    ssl_mode: ${POSTGRES_SSL_MODE}
	//trace:
	//  endpoint: ${JAEGER_ENDPOINT}
	//#sentry:
	//#  dsn: ${SENTRY_DSN}
	//`
)

func ArgDocker() error {
	services := &settingsPackage.Services{}
	if err := services.Load(); err != nil {
		return err
	}
	fmt.Println("\n--- array")
	for _, service := range services.Services {
		fmt.Printf("      - %s\n", service.NameVersion)
	}
	fmt.Println("\n--- docker-compose")
	tempService, err := template.New("service").Parse(service)
	if err != nil {
		return err
	}
	for p, service := range services.Services {
		service.Port = int16(p + 9001)
		service.HttpPort = int16(p + 8001)
		if err := tempService.Execute(os.Stdout, service); err != nil {
			return err
		}
	}
	//fmt.Println("\n--- config")
	//tempConfig, err := template.New("config").Parse(config)
	//if err != nil {
	//	return err
	//}
	//for _, service := range services.Services {
	//	if err := tempConfig.Execute(os.Stdout, service); err != nil {
	//		return err
	//	}
	//}
	return nil
}
