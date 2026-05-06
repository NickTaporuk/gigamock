package main

import (
	"context"
	"log"

	"github.com/NickTaporuk/gigamock/src/app"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func newRootCommand() *cobra.Command {
	cfg, err := app.DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd := &cobra.Command{
		Use:   "gigamock",
		Short: "Generic mock server for HTTP, Kafka, and future dynamic gRPC mocks",
		Long: `Gigamock serves mock responses described in YAML or JSON files.

It indexes a mock configuration directory, starts the mock HTTP server, and
exposes a built-in control UI for switching active scenarios at runtime.`,
		Example: `  gigamock --dir-path ./examples/rest
  gigamock --dir-path ./examples/rest --dir-path ./examples/graphql --dir-path ./examples/grpc
  gigamock --server-ip 127.0.0.1 --server-port :7777 --grpc-server-port :7778
  gigamock --logger-level INFO --logger-pretty-print`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			inst := app.NewApp(ctx)
			defer func() {
				cancel()
				if err := inst.Stop(); err != nil {
					log.Fatal(err)
				}
			}()

			return inst.RunWithConfig(cfg)
		},
	}

	rootCmd.Flags().StringVar(&cfg.ServerIP, "server-ip", cfg.ServerIP, "server IP address to bind")
	rootCmd.Flags().StringVar(&cfg.ServerPort, "server-port", cfg.ServerPort, "server port to listen on")
	rootCmd.Flags().StringVar(&cfg.GRPCServerPort, "grpc-server-port", cfg.GRPCServerPort, "gRPC server port to listen on")
	rootCmd.Flags().StringArrayVar(&cfg.DirPaths, "dir-path", cfg.DirPaths, "mock configuration directory with YAML or JSON files; can be used multiple times")
	rootCmd.Flags().StringVar(&cfg.LoggerLevel, "logger-level", cfg.LoggerLevel, "logger level: DEBUG, INFO, WARN, ERROR")
	rootCmd.Flags().BoolVar(&cfg.LoggerPrettyPrint, "logger-pretty-print", cfg.LoggerPrettyPrint, "enable human-readable pretty log output")

	return rootCmd
}
