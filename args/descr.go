package args

import (
	"fmt"
)

func ShowDescription() {
	fmt.Println("This util for services deployment automation")
	fmt.Println("-upv, -uvp, -vup - up the version")
	fmt.Println("-env - envoy some settings")
	fmt.Println("-mak - make commands")
	fmt.Println("-dcr - docker make commands")
	fmt.Println("-prt service - make proto rpc and messages with service")
	fmt.Println("-rep service - implements services")
	fmt.Println("-sql service - implements getReply etc")
	fmt.Println("-srv service - get commands of msg, enm, req (pinia store)")
	fmt.Println("-msg message - make json get & set message")
	fmt.Println("-req message - make act request")
	fmt.Println("-enm enum - make json case & array enum")
	fmt.Println("-str service - make pinia actions of service")
}
