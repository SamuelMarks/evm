package commands

import (
	"fmt"
	"log/syslog"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"

	"github.com/Fantom-foundation/go-evm/src/engine"
)

//AddRunFlags adds flags to the Run command
func AddRunFlags(cmd *cobra.Command) {
	//Lachesis Socket
	cmd.Flags().String("proxy", config.ProxyAddr, "IP:PORT of Lachesis proxy")
	if runtime.GOOS != "windows" {
		cmd.Flags().Bool("syslog", config.Syslog, "IP:PORT of Lachesis proxy")
	}
	viper.BindPFlags(cmd.Flags())
}

// NewRunCmd returns the command that allows the CLI to start a node.
func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the evm node",
		RunE:  run,
	}

	AddRunFlags(cmd)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {

	if config.Syslog {
		hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
		if err == nil {
			logger.Hooks.Add(hook)
		}
	}
	engine, err := engine.NewSocketEngine(*config, logger)
	//engine, err := engine.NewInmemEngine(*config, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
