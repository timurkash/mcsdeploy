package args

import "fmt"

func ShowDescription() {
	fmt.Println("This util for services deployment automation")
	fmt.Println("-env - envoy")
	fmt.Println("-upv, -uvp, -vup - for upping the version")
	fmt.Println("Util run on local dir with services.yaml")
}
