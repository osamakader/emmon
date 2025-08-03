package main

import (
	"fmt"
	"os"

	"emmon/monitor"
	"emmon/terminal"
	"emmon/web"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	log     = logrus.New()
)

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "emmon",
	Short: "Embedded Linux System Monitor",
	Long: `A lightweight system monitor for embedded Linux devices.
Features:
- Real-time CPU, RAM, temperature, and disk monitoring
- GPIO status monitoring
- Web UI with WebSocket updates
- Terminal UI with tcell
- Lightweight and optimized for embedded systems`,
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web interface",
	Long:  `Start the embedded monitor with web UI and WebSocket support`,
	Run: func(cmd *cobra.Command, args []string) {
		port := viper.GetString("web.port")
		log.Infof("Starting web interface on port %s", port)
		startWebInterface(port)
	},
}

var terminalCmd = &cobra.Command{
	Use:   "terminal",
	Short: "Start the terminal interface",
	Long:  `Start the embedded monitor with terminal UI using tcell`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting terminal interface")
		startTerminalInterface()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.emmon.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")

	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))

	// Web command flags
	webCmd.Flags().String("port", "8080", "port for web interface")
	viper.BindPFlag("web.port", webCmd.Flags().Lookup("port"))

	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(terminalCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".emmon")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	// Set up logging
	logLevel := viper.GetString("log.level")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Warnf("Invalid log level %s, using info", logLevel)
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// startWebInterface starts the web interface
func startWebInterface(port string) {
	monitor := monitor.NewSystemMonitor(log)
	server := web.NewWebServer(port, log, monitor)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}

// startTerminalInterface starts the terminal interface
func startTerminalInterface() {
	monitor := monitor.NewSystemMonitor(log)
	ui := terminal.NewTerminalUI(monitor, log)

	if err := ui.Start(); err != nil {
		log.Fatalf("Failed to start terminal UI: %v", err)
	}
}
