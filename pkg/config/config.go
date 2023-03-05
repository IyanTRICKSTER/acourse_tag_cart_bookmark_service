package config

import (
	"acourse_tag_cart_bookmark_service/pkg/contracts"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	App      map[string]string
	Database map[string]string
}

func Construct(envpath string) contracts.Config {

	//Load .Env
	err := godotenv.Load(envpath)
	if err != nil {
		panic(err)
	}

	c := Config{}

	c.App = map[string]string{}
	c.App["PORT"] = os.Getenv("APP_PORT")
	c.App["RPC_TARGET_HOST"] = os.Getenv("RPC_TARGET_HOST")
	c.App["RPC_TARGET_PORT"] = os.Getenv("RPC_TARGET_PORT")

	c.Database = map[string]string{}
	c.Database["USERNAME"] = os.Getenv("DB_USERNAME")
	c.Database["PASSWORD"] = os.Getenv("DB_PASSWORD")
	c.Database["HOST"] = os.Getenv("DB_HOST")
	c.Database["PORT"] = os.Getenv("DB_PORT")
	c.Database["NAME"] = os.Getenv("DB_NAME")
	c.Database["COLLECTION_BOOKMARKS"] = os.Getenv("DB_COLLECTION_BOOKMARKS")
	c.Database["COLLECTION_TAGS"] = os.Getenv("DB_COLLECTION_TAGS")
	c.Database["COLLECTION_CARTS"] = os.Getenv("DB_COLLECTION_CARTS")

	return &c
}

func (c *Config) GetDBConfig() map[string]string {
	return c.Database
}

func (c Config) GetAppConfig() map[string]string {
	return c.App
}
