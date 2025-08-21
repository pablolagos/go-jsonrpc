package examples

import (
	"fmt"
	"time"

	"github.com/pablolagos/go-jsonrpc/jclient"
)

func main() {
	client := jclient.NewTCPClient("127.0.0.1:9000", 5*time.Second)

	var result map[string]interface{}
	err := client.Call("ping", map[string]string{"msg": "hello"}, &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", result)
}
