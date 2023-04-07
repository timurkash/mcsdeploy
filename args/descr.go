package args

import (
	"fmt"
)

func ShowDescription() {
	fmt.Println("This util for services deployment automation")
	fmt.Println("-env - envoy some settings")
	fmt.Println("-doc - docker make commands")
	//fmt.Println("-upv, -uvp, -vup - for upping the version (deprecated)")
	fmt.Println("-prt - make proto rpc and messages")
	fmt.Println("Util run on local dir with services.yaml")
}
