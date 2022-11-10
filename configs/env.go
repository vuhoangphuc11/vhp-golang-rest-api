package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGOURI")
}

//func GetConfig(str string) interface{} {
//	viper.SetConfigFile(".env")
//	err := viper.ReadInConfig()
//	if err != nil {
//		viper.SetConfigFile("config.yaml")
//	}
//	viper.ReadInConfig()
//
//	return viper.Get(str)
//}
