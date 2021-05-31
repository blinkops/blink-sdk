package main

import (
	"fmt"
	_ "github.com/blinkops/blink-sdk/plugin"
	_ "github.com/blinkops/blink-sdk/plugin/actions"
	_ "github.com/blinkops/blink-sdk/plugin/config"
	_ "github.com/blinkops/blink-sdk/plugin/connections"
	_ "github.com/blinkops/blink-sdk/plugin/description"
	_ "github.com/blinkops/blink-sdk/plugin/proto"
	_ "github.com/blinkops/blink-sdk/plugin/server"
)

func main() {
	fmt.Println("This is Blink Plugins SDK main")
}
