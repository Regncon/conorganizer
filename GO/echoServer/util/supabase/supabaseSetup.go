package supabaseSetup

import (
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
	"regncon.no/htmx/echoServer/util"
)

var Client *supabase.Client

func envVariable(key string) string {
	if env := os.Getenv(key); env == "" {
		return util.GoDotEnvVariable(key)
	}
	return os.Getenv(key)
}

func initiateClient(supabaseClient *supabase.Client) {
	Client = supabaseClient
	fmt.Println("client intiated")
}

func Init() {
	API_URL := envVariable("SUPABASE_API_URL")
	API_KEY := envVariable("SUPABASE_API_KEY")
	client, err := supabase.NewClient(API_URL, API_KEY, nil)
	initiateClient(client)
	if err != nil {
		fmt.Println("err", err)
	}
}
