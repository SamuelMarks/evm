package engine

import (
	"fmt"
	//"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Fantom-foundation/go-evm/src/config"
	"github.com/Fantom-foundation/go-evm/src/service"
	"github.com/Fantom-foundation/go-evm/src/state"
	"github.com/SamuelMarks/dag1/src/crypto"
	"github.com/SamuelMarks/dag1/src/net"
	"github.com/SamuelMarks/dag1/src/node"
	"github.com/SamuelMarks/dag1/src/peers"
	"github.com/SamuelMarks/dag1/src/poset"
	serv "github.com/SamuelMarks/dag1/src/service"
)

type InmemEngine struct {
	ethService *service.Service
	ethState   *state.State
	node       *node.Node
	service    *serv.Service
}

func NewInmemEngine(config config.Config, logger *logrus.Logger) (*InmemEngine, error) {
	submitCh := make(chan []byte)

	state, err := state.NewState(logger,
		config.Eth.DbFile,
		config.Eth.Cache)
	if err != nil {
		return nil, err
	}

	service := service.NewService(config.Eth.Genesis,
		config.Eth.Keystore,
		config.Eth.EthAPIAddr,
		config.Eth.PwdFile,
		state,
		submitCh,
		logger)

	appProxy := NewInmemProxy(state, service, submitCh, logger)

	//------------------------------------------------------------------------------

	// Create the PEM key
	pemKey := crypto.NewPemKey(config.DAG1.DataDir)

	// Try a read
	key, err := pemKey.ReadKey()
	if err != nil {
		return nil, err
	}

	// Create the peer store
	peerStore := peers.NewJSONPeers(config.DAG1.DataDir)
	// Try a read
	participants, err := peerStore.Peers()
	if err != nil {
		return nil, err
	}

	// There should be at least two peers
	if participants.Len() < 2 {
		return nil, fmt.Errorf("peers.json should define at least two peers")
	}

	pmap := participants

	//Find the ID of this node
	nodePub := fmt.Sprintf("0x%X", crypto.FromECDSAPub(&key.PublicKey))
	n, ok := pmap.ByPubKey[nodePub]

	if !ok {
		return nil, fmt.Errorf("cannot find self pubkey in peers.json")
	}

	nodeID := n.ID

	logger.WithFields(logrus.Fields{
		"pmap": pmap,
		"id":   nodeID,
	}).Debug("Participants")

	conf := node.NewConfig(
		time.Duration(config.DAG1.Heartbeat)*time.Millisecond,
		time.Duration(config.DAG1.TCPTimeout)*time.Millisecond,
		config.DAG1.CacheSize,
		config.DAG1.SyncLimit,
		logger)

	//Instantiate the Store (inmem or badger)
	var store poset.Store
	//var needBootstrap bool
	/* TODO inmem only for now */
	/*switch conf.StoreType {
	case "inmem":*/
	store = poset.NewInmemStore(pmap, conf.CacheSize)
	/*case "badger":
		//If the file already exists, load and bootstrap the store using the file
		if _, err := os.Stat(conf.StorePath); err == nil {
			logger.Debug("loading badger store from existing database")
			store, err = poset.LoadBadgerStore(conf.CacheSize, conf.StorePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load BadgerStore from existing file: %s", err)
			}
			needBootstrap = true
		} else {
			//Otherwise create a new one
			logger.Debug("creating new badger store from fresh database")
			store, err = poset.NewBadgerStore(pmap, conf.CacheSize, conf.StorePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create new BadgerStore: %s", err)
			}
		}
	default:
		return nil, fmt.Errorf("Invalid StoreType: %s", conf.StoreType)
	}*/

	trans, err := net.NewTCPTransport(
		config.DAG1.BindAddr, nil, 2, conf.TCPTimeout, logger)
	if err != nil {
		return nil, fmt.Errorf("creating TCP Transport: %s", err)
	}

	node := node.NewNode(conf, nodeID, key, participants, store, trans, appProxy)
	if err := node.Init(); err != nil {
		return nil, fmt.Errorf("initializing node: %s", err)
	}

	lserv := serv.NewService(config.DAG1.ServiceAddr, node, logger)

	return &InmemEngine{
		ethState:   state,
		ethService: service,
		node:       node,
		service:    lserv,
	}, nil

}

/*******************************************************************************
Implement Engine interface
*******************************************************************************/

func (i *InmemEngine) Run() error {

	//ETH API service
	go i.ethService.Run()

	//DAG1 API service
	go i.service.Serve()

	i.node.Run(true)

	return nil
}
