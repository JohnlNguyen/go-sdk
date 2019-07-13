package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"go-sdk/logger"
	"go-sdk/secrets"
	"go-sdk/yaml"
)

type myConfig struct {
	Environment string `yaml:"environment" env:"SERVICE_ENV" secret:"environment"`
	Version     string `yaml:"version" env:"SERVICE_VERSION"`
	Secret      string `yaml:"secret" secret:"secret,base64"`
}

func main() {
	log := logger.All()
	client := secrets.Must(secrets.NewFromEnv()).WithLogger(log)

	keyPath := "configTest"
	err := client.WriteInto(context.TODO(), keyPath, myConfig{
		Environment: "test",
		Version:     "wont be in output",
		Secret:      base64.StdEncoding.EncodeToString([]byte("a super secure one")),
	})
	if err != nil {
		log.SyncFatalExit(err)
	}

	var cfg myConfig
	err = client.ReadInto(context.TODO(), keyPath, &cfg)
	if err != nil {
		log.SyncFatalExit(err)
	}

	contents, err := yaml.Marshal(cfg)
	if err != nil {
		log.SyncFatalExit(err)
	}
	fmt.Printf("%v\n", string(contents))
}
