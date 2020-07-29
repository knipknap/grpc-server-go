package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Debug bool `env:"DEBUG" envDefault:"false"`
}

var Data config

func init() {
	Data = config{}
	if err := env.Parse(&Data); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
