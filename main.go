package main

import (
	"github.com/fulldump/goconfig"

	"go-solidary/api"
	"go-solidary/config"
)

func main() {
	c := config.Config{}
	goconfig.Read(&c)

	err := api.UpServer(&c)
	if err != nil {
		panic(err)
	}
}
