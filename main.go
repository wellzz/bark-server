package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"os"
	"strconv"
)

var listenAddr net.IP
var listenPort int
var debug bool
var dev bool
var dataDir string

var rootCmd = &cobra.Command{
	Use:   "bark-server",
	Short: "Bark Server",
	Long: `
Bark Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
			logrus.Info("rewrite listen port:", val)
			listenPort, _ = strconv.Atoi(val)
		}
		runBarkServer()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().IPVarP(&listenAddr, "listen", "l", net.ParseIP("0.0.0.0"), "server listen address")
	rootCmd.PersistentFlags().IntVarP(&listenPort, "port", "p", 8080, "server listen port")
	rootCmd.PersistentFlags().StringVarP(&dataDir, "data", "d", "/data", "data dir")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().BoolVar(&dev, "dev", false, "dev mode")
}
