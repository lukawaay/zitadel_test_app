package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lukawaay/zitadel_test_app/internal/server"
	"github.com/lukawaay/zitadel_test_app/internal/server/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := config.Load(config.Source {
		InstanceDomain: os.Getenv("ZITADEL_TEST_APP_INSTANCE_DOMAIN"),
		InstancePort: os.Getenv("ZITADEL_TEST_APP_INSTANCE_PORT"),
		InstanceProtocol: os.Getenv("ZITADEL_TEST_APP_INSTANCE_PROTOCOL"),
		Key: os.Getenv("ZITADEL_TEST_APP_KEY"),
		ClientID: os.Getenv("ZITADEL_TEST_APP_CLIENT_ID"),
		RedirectURI: os.Getenv("ZITADEL_TEST_APP_REDIRECT_URI"),
		Port: os.Getenv("ZITADEL_TEST_APP_PORT"),
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load the config: %s", err))
		os.Exit(1)
	}

	if err := server.Start(config, logger); err != nil {
		logger.Error(fmt.Sprintf("Failed to start the server: %s", err))
		os.Exit(1)
	}
}
