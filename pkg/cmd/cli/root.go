package cli

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Version: "1.2.3",
	Use:     "blitz-cli",
	Short:   "blitz-cli is the CLI for blitzd",
	Long: `An easy way to access data from blitzd.
                More info at: https://github.com/frennkie/blitzd`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		log.WithFields(log.Fields{"rpcHostPort": viper.GetString("rpcHostPort")}).Debug("cli config")
	},
}

func Init() {

	rootCmd.PersistentFlags().StringVarP(&config.BlitzdDir, "dir", "D",
		config.DefaultBlitzdDir, "blitzd home directory")

	rootCmd.PersistentFlags().StringVarP(&config.RpcHostPort,
		"rpcHostPort", "H", fmt.Sprintf("localhost:%d", config.DefaultRPCPort),
		"Host and Port to connect to")
	_ = viper.BindPFlag("rpcHostPort", rootCmd.PersistentFlags().Lookup("rpcHostPort"))

	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v",
		false, "print debug log messages")

	rootCmd.PersistentFlags().BoolVarP(&config.Trace, "trace", "t",
		false, "print all (also debug and trace) log messages")
	_ = rootCmd.PersistentFlags().MarkHidden("trace")

	rootCmd.AddCommand(cmdTimes)
	rootCmd.AddCommand(cmdEcho)
	rootCmd.AddCommand(cmdHello)
	rootCmd.AddCommand(cmdHelloWorld)

	rootCmd.AddCommand(cmdGet)
	cmdGet.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output as JSON")
	cmdGet.Flags().BoolVarP(&formattedFlag, "formatted", "f", false, "Output as formatted value")

	cmdGet.AddCommand(cmdGetAll)
	cmdGet.AddCommand(cmdGetFoo5)

	rootCmd.AddCommand(cmdFoo5)

	_ = rootCmd.Execute()

}
