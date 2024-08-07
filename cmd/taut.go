package main

import (
	"github.com/ntietz/taut-chat/pkg/web"
)

func main() {
	s := web.CreateServer()
	s.Logger.Fatal(s.Start(":8080"))
}
