package main

import (
	"fmt"

	"regncon.no/htmx/echoServer"
	supabaseSetup "regncon.no/htmx/echoServer/util/supabase"
)

func init() {
	supabaseSetup.Init()
}

func main() {
	fmt.Println("Hello, World!")
	echoServer.EchoServer()
}
