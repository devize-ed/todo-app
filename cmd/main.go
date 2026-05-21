package cmd

import (
	"fmt"

	"github.com/devize-ed/todo-app/internal/config"
	"github.com/devize-ed/todo-app/internal/logger"
)

func main() {

}

func run() error {
	cfg, err := config.Load(config.GetFilePath())
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if _, err := logger.Initialize(cfg.LogLevel); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.SafeSync()

	return nil
}
