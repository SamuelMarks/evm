package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Fantom-foundation/go-evm/src/consensus/dag1"
	"github.com/Fantom-foundation/go-evm/src/engine"
)

//AddDAG1Flags adds flags to the DAG1 command
func AddDAG1Flags(cmd *cobra.Command) {
	cmd.Flags().String("dag1.datadir", config.DAG1.DataDir, "Directory containing priv_key.pem and peers.json files")
	cmd.Flags().String("dag1.listen", config.DAG1.BindAddr, "IP:PORT of DAG1 node")
	cmd.Flags().String("dag1.service-listen", config.DAG1.ServiceAddr, "IP:PORT of DAG1 HTTP API service")
	cmd.Flags().Duration("dag1.heartbeat", config.DAG1.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Duration("dag1.timeout", config.DAG1.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("dag1.cache-size", config.DAG1.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int64("dag1.sync-limit", config.DAG1.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("dag1.max-pool", config.DAG1.MaxPool, "Max number of pool connections")
	cmd.Flags().Bool("dag1.store", config.DAG1.Store, "use persistent store")
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic("Unable to bind viper flags")
	}
}

//NewDAG1Cmd returns the command that starts EVM-Lite with DAG1 consensus
func NewDAG1Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dag1",
		Short: "Run the evm node with DAG1 consensus",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			config.SetDataDir(config.BaseConfig.DataDir)

			logger.WithFields(logrus.Fields{
				"DAG1": config.DAG1,
			}).Debug("Config")

			return nil
		},
		RunE: runDAG1,
	}

	AddDAG1Flags(cmd)

	return cmd
}

func runDAG1(cmd *cobra.Command, args []string) error {
	dag1Consensus := dag1.NewInmemDAG1(config.DAG1, logger)
	consensusEngine, err := engine.NewConsensusEngine(*config, dag1Consensus, logger)
	if err != nil {
		return fmt.Errorf("error building Engine: %s", err)
	}

	if err := consensusEngine.Run(); err != nil {
		return fmt.Errorf("error running Engine: %s", err)
	}

	return nil
}
