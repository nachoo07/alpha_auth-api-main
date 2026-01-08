package main

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/melisource/fury_go-core/pkg/log"

// 	"github.com/melisource/fury_shp-travel-management-api/internal/config"
// 	"github.com/melisource/fury_shp-travel-management-api/internal/migration"
// )

// const (
// 	upCommand         = "UP"
// 	downCommand       = "DOWN"
// 	usageMessage      = "Usage: migration.go (UP|DOWN)"
// 	unexpectedMessage = "unexpected err"
// 	expectedArgsLen   = 2
// )

// func main() {
// 	lvl := log.NewAtomicLevel()
// 	logger := log.NewProductionLogger(&lvl)

// 	if !config.IsLocalScope() {
// 		logger.Error("migration can only be executed in local scope")
// 		os.Exit(1)
// 	}

// 	cmd, err := parseArgs()
// 	if err != nil {
// 		logger.Error(unexpectedMessage, log.Err(err), log.String("usage", usageMessage))
// 		os.Exit(1)
// 	}

// 	cfg, err := config.NewConfig()
// 	if err != nil {
// 		logger.Error(unexpectedMessage, log.Err(err))
// 		os.Exit(1)
// 	}

// 	if err = executeMigration(cfg, cmd); err != nil {
// 		logger.Error(unexpectedMessage, log.Err(err))
// 		os.Exit(1)
// 	}

// 	logger.Info("Migration successfully processed")
// }

// func parseArgs() (string, error) {
// 	if len(os.Args) != expectedArgsLen {
// 		return "", errors.New("invalid number of arguments")
// 	}

// 	switch arg := strings.ToUpper(os.Args[1]); arg {
// 	case upCommand:
// 		return upCommand, nil
// 	case downCommand:
// 		return downCommand, nil
// 	default:
// 		return "", fmt.Errorf("invalid argument: %v", arg)
// 	}
// }

// func executeMigration(cfg config.Configuration, cmd string) error {
// 	migrator := migration.NewClient(cfg)

// 	if downCommand == cmd {
// 		return migrator.Down()
// 	}

// 	return migrator.Up()
// }
