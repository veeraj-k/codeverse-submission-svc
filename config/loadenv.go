package config

import (
	"github.com/joho/godotenv"
)

func Loadenv() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

}
