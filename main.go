package main

import (
	"fmt"
	"os"

	eirinix "github.com/SUSE/eirinix"
	helloworld "github.com/drnic/eirinix-helloworld/hello"
)

func main() {
	fmt.Println("Running drnic/eirinix-helloworld...")
	options := eirinix.ManagerOptions{
		Namespace:           os.Getenv("POD_NAMESPACE"),
		Host:                "0.0.0.0",
		Port:                4545,
		ServiceName:         os.Getenv("WEBHOOK_SERVICE_NAME"),
		WebhookNamespace:    os.Getenv("WEBHOOK_NAMESPACE"),
		OperatorFingerprint: "eirini-x-drnic-helloworld",
	}
	fmt.Printf("--> %#v\n", options)
	x := eirinix.NewManager(options)
	x.AddExtension(&helloworld.Extension{})
	x.Start()
}
