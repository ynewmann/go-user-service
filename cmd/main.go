package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"go-user-service/src/controllers"
	"go-user-service/src/handlers"
	"go-user-service/src/repository"
	"go-user-service/src/repository/postgres"
	"go-user-service/src/server"
)

var ErrInvalidConfigFileName = errors.New("invalid config file name")

const (
	config = "config"
	cfg    = "c"
)

type Config struct {
	Server   server.Config     `yaml:"server"`
	Database repository.Config `yaml:"database"`
}

var rootCmd = &cobra.Command{
	Use:   "userservice",
	Short: "user service",
	Long:  "user service",
	RunE:  runRootCmd,
}

func init() {
	rootCmd.Flags().StringP(config, cfg, "", "path to config file")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	err = cmd.Flags().Parse(args)
	if err != nil {
		logger.Error("failed to parse flags", zap.Error(err))
		return err
	}

	configFile, err := cmd.Flags().GetString(config)
	if err != nil {
		logger.Error("failed to get config path", zap.Error(err))
		return err
	}

	splitPath := strings.Split(configFile, "/")
	fullName := splitPath[len(splitPath)-1]
	splitName := strings.Split(fullName, ".")
	if len(splitName) != 2 {
		logger.Error("cannot parse config", zap.String("file name", fullName))
		return ErrInvalidConfigFileName
	}

	viper.SetConfigName(splitName[0])
	viper.SetConfigType(splitName[1])
	viper.AddConfigPath(strings.Join(splitPath[:len(splitPath)-1], "/"))

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("cannot read config", zap.Error(err))
		return err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.Error("unable to decode config file into struct", zap.Error(err))
		return err
	}

	repo, err := postgres.NewRepository(config.Database)
	if err != nil {
		logger.Error("cannot create repo", zap.Error(err))
		return err
	}

	userController := controllers.New(repo)
	userHandler := handlers.New(userController)
	microservice := server.New(config.Server, logger, userHandler)

	err = microservice.Start()
	if err != nil {
		logger.Error("while running server", zap.Error(err))
		return err
	}

	return nil
}
