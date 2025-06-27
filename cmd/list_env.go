package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/intezya/auth_service/internal/adapters/config"
)

func main() {
	var cfg config.Config
	desc, _ := cleanenv.GetDescription(&cfg, nil)
	fmt.Println(desc)
}
