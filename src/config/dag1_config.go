package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	_dag1 "github.com/SamuelMarks/dag1/src/dag1"
)

const (
	defaultNodeAddr        = ":1337"
	defaultDAG1APIAddr = ":8000"
	defaultHeartbeat       = 500 * time.Millisecond
	defaultTCPTimeout      = 1000 * time.Millisecond
	defaultCacheSize       = 50000
	defaultSyncLimit       = 1000
	defaultMaxPool         = 2
)

var (
	defaultDAG1Dir = fmt.Sprintf("%s/dag1", DefaultDataDir)
	defaultPeersFile   = fmt.Sprintf("%s/peers.json", defaultDAG1Dir)
)

// DAG1Config contains the configuration of a DAG1 node
type DAG1Config struct {
	// Directory containing priv_key.pem and peers.json files
	DataDir string `mapstructure:"datadir"`

	// Address of DAG1 node (where it talks to other DAG1 nodes)
	BindAddr string `mapstructure:"listen"`

	// DAG1 HTTP API address
	ServiceAddr string `mapstructure:"service-listen"`

	// Gossip heartbeat
	Heartbeat time.Duration `mapstructure:"heartbeat"`

	// TCP timeout
	TCPTimeout time.Duration `mapstructure:"timeout"`

	// Max number of items in caches
	CacheSize int `mapstructure:"cache-size"`

	// Max number of Event in SyncResponse
	SyncLimit int64 `mapstructure:"sync-limit"`

	// Max number of connections in net pool
	MaxPool int `mapstructure:"max-pool"`

	// Database type; badger or inmeum
	Store bool `mapstructure:"store"`
}

// DefaultDAG1Config returns the default configuration for a DAG1 node
func DefaultDAG1Config() *DAG1Config {
	return &DAG1Config{
		DataDir:     defaultDAG1Dir,
		BindAddr:    defaultNodeAddr,
		ServiceAddr: defaultDAG1APIAddr,
		Heartbeat:   defaultHeartbeat,
		TCPTimeout:  defaultTCPTimeout,
		CacheSize:   defaultCacheSize,
		SyncLimit:   defaultSyncLimit,
		MaxPool:     defaultMaxPool,
	}
}

// SetDataDir updates the dag1 configuration directories if they were set to
// to default values.
func (c *DAG1Config) SetDataDir(datadir string) {
	if c.DataDir == defaultDAG1Dir {
		c.DataDir = datadir
	}
}

// ToRealDAG1Config converts an evm/src/config.DAG1Config to a
// dag1/src/dag1.DAG1Config as used by DAG1
func (c *DAG1Config) ToRealDAG1Config(logger *logrus.Logger) *_dag1.DAG1Config {
	dag1Config := _dag1.NewDefaultConfig()
	dag1Config.DataDir = c.DataDir
	dag1Config.BindAddr = c.BindAddr
	dag1Config.ServiceAddr = c.ServiceAddr
	dag1Config.MaxPool = c.MaxPool
	dag1Config.Store = c.Store
	dag1Config.Logger = logger
	dag1Config.NodeConfig.HeartbeatTimeout = c.Heartbeat
	dag1Config.NodeConfig.TCPTimeout = c.TCPTimeout
	dag1Config.NodeConfig.CacheSize = c.CacheSize
	dag1Config.NodeConfig.SyncLimit = c.SyncLimit
	dag1Config.NodeConfig.Logger = logger
	return dag1Config
}
