package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	// Base
	defaultLogLevel = "debug"
	DefaultDataDir  = defaultHomeDir()
)

// Config contains de configuration for an EVM-Lite node
type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	// Options for EVM and State
	Eth *EthConfig `mapstructure:"eth"`

	// Options for DAG1 consensus
	DAG1 *DAG1Config `mapstructure:"dag1"`

	// Options for Raft consensus
	Raft *RaftConfig `mapstructure:"raft"`

	ProxyAddr  string `mapstructure:"proxy"`
	ClientAddr string `mapstructure:"client-connect"`
	Standalone bool   `mapstructure:"standalone"`
	Pidfile    string `mapstructure:"pidfile"`
}

// DefaultConfig returns the default configuration for an EVM-Lite node
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Eth:        DefaultEthConfig(),
		DAG1:   DefaultDAG1Config(),
		Raft:       DefaultRaftConfig(),
		ProxyAddr:  "127.0.0.1:1338",
		ClientAddr: "127.0.0.1:1339",
		Pidfile:    filepath.Join(os.TempDir(), "go-evm.pid"),
	}
}

// SetDataDir updates the root data directory as well as the various lower config
// for eth and consensus
func (c *Config) SetDataDir(datadir string) {
	c.BaseConfig.DataDir = datadir
	if c.Eth != nil {
		c.Eth.SetDataDir(fmt.Sprintf("%s/eth", datadir))
	}
	if c.DAG1 != nil {
		c.DAG1.SetDataDir(fmt.Sprintf("%s/dag1", datadir))
	}
	if c.Raft != nil {
		c.Raft.SetDataDir(fmt.Sprintf("%s/raft", datadir))
	}
}

/*******************************************************************************
BASE CONFIG
*******************************************************************************/

// BaseConfig contains the top level configuration for an EVM-DAG1 node
type BaseConfig struct {
	// Top-level directory of evm-dag1 data
	DataDir string `mapstructure:"datadir"`

	// Debug, info, warn, error, fatal, panic
	LogLevel string `mapstructure:"log"`
}

// DefaultBaseConfig returns the default top-level configuration for EVM-DAG1
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		DataDir:  DefaultDataDir,
		LogLevel: defaultLogLevel,
	}
}

/*******************************************************************************
FILE HELPERS
*******************************************************************************/

func defaultHomeDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "EVM")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "EVM")
		} else {
			return filepath.Join(home, ".evm")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
