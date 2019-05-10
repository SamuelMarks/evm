package dag1

import (
	"github.com/sirupsen/logrus"

	"github.com/Fantom-foundation/go-evm/src/config"
	"github.com/Fantom-foundation/go-evm/src/service"
	"github.com/Fantom-foundation/go-evm/src/state"
	_dag1 "github.com/SamuelMarks/dag1/src/dag1"
)

// InmemDAG1 implements the Consensus interface.
// It uses an inmemory DAG1 node.
type InmemDAG1 struct {
	config     *config.DAG1Config
	dag1   *_dag1.DAG1
	ethService *service.Service
	ethState   *state.State
	logger     *logrus.Logger
}

// NewInmemDAG1 instantiates a new InmemDAG1 consensus system
func NewInmemDAG1(config *config.DAG1Config, logger *logrus.Logger) *InmemDAG1 {
	return &InmemDAG1{
		config: config,
		logger: logger,
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

// Init instantiates a DAG1 inmemory node
func (b *InmemDAG1) Init(state *state.State, service *service.Service) error {

	b.logger.Debug("INIT")

	b.ethState = state
	b.ethService = service

	realConfig := b.config.ToRealDAG1Config(b.logger)
	realConfig.Proxy = NewInmemProxy(state, service, service.GetSubmitCh(), b.logger)

	dag1 := _dag1.NewDAG1(realConfig)
	err := dag1.Init()
	if err != nil {
		return err
	}
	b.dag1 = dag1

	return nil
}

// Run starts the DAG1 node
func (b *InmemDAG1) Run() error {
	b.dag1.Run()
	return nil
}

// Info returns DAG1 stats
func (b *InmemDAG1) Info() (map[string]string, error) {
	info := b.dag1.Node.GetStats()
	info["type"] = "dag1"
	return info, nil
}
