package main

import (
	"fmt"
	_ "github.com/blinkops/plugin-sdk/plugin"
	_ "github.com/blinkops/plugin-sdk/plugin/actions"
	_ "github.com/blinkops/plugin-sdk/plugin/config"
	_ "github.com/blinkops/plugin-sdk/plugin/connections"
	_ "github.com/blinkops/plugin-sdk/plugin/description"
	_ "github.com/blinkops/plugin-sdk/plugin/proto"
	_ "github.com/blinkops/plugin-sdk/plugin/server"
)

func main() {
	fmt.Println("This is Blink Plugins SDK main")
}
