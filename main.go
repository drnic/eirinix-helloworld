package main

import (
	"fmt"
	"os"

	eirinix "github.com/SUSE/eirinix"
	helloworld "github.com/drnic/eirinix-helloworld/hello"
)

func main() {
	fmt.Println("Running drnic/eirinix-helloworld...")
	x := eirinix.NewManager(
		eirinix.ManagerOptions{
			Namespace:           os.Getenv("NAMESPACE"),
			WebhookNamespace:    os.Getenv("NAMESPACE"),
			Host:                "0.0.0.0",
			Port:                4545,
			ServiceName:         os.Getenv("WEBHOOK_SERVICE_NAME"),
			OperatorFingerprint: "eirini-x-drnic-helloworld",
		})
	x.AddExtension(&helloworld.Extension{})
	x.Start()
}
