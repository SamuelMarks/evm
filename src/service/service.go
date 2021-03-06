package service

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/Fantom-foundation/go-evm/src/common"
	"github.com/Fantom-foundation/go-evm/src/config"
	"github.com/Fantom-foundation/go-evm/src/state"
)

var defaultGas = hexutil.Uint64(90000)

type infoCallback func() (map[string]string, error)

type Service struct {
	sync.Mutex
	chainConfig *params.ChainConfig
	state       *state.State
	submitCh    chan []byte
	genesisFile string
	keystoreDir string
	apiAddr     string
	keyStore    *keystore.KeyStore
	am          *accounts.Manager
	pwdFile     string
	logger      *logrus.Logger

	rpcConfig *node.Config
	rpcServer *RpcServer

	//XXX
	getInfo infoCallback
}

func NewService(genesisFile, keystoreDir, apiAddr, pwdFile string,
	state *state.State,
	submitCh chan []byte,
	logger *logrus.Logger) *Service {
	// TODO: replace DefaultRpcConfig with custom
	rpcConfig := &config.DefaultRpcConfig
	// TODO: replace ChainConfig with custom
	chainConfig := &params.ChainConfig{
		ChainID: big.NewInt(666),
	}

	s := &Service{
		chainConfig: chainConfig,
		genesisFile: genesisFile,
		keystoreDir: keystoreDir,
		apiAddr:     apiAddr,
		pwdFile:     pwdFile,
		state:       state,
		submitCh:    submitCh,
		logger:      logger,
		// TODO: no-default rpcConfig required
		rpcConfig: rpcConfig,
	}
	var err error
	s.rpcServer, err = NewRpcServer(rpcConfig, s)
	if err != nil {
		panic(err)
	}
	err = s.rpcServer.Register(NewWeb3AccountServiceConstructor(s))
	if err != nil {
		panic(err)
	}
	err = s.rpcServer.Register(NewWeb3ChainServiceConstructor(s))
	if err != nil {
		panic(err)
	}

	return s
}

func (m *Service) Run() {
	m.checkErr(m.makeKeyStore())
	m.checkErr(m.unlockAccounts())
	m.checkErr(m.createGenesisAccounts())

	m.logger.Info("serving web3-api ...")
	if err := m.rpcServer.Start(); err != nil {
		panic(err)
	}
	defer (func() {
		if err := m.rpcServer.Stop(); err != nil {
			panic(err)
		}
	})()

	m.logger.Info("serving api ...")
	m.serveAPI()
}

//XXX
func (m *Service) GetSubmitCh() chan []byte {
	return m.submitCh
}

//XXX
func (m *Service) SetInfoCallback(f infoCallback) {
	m.getInfo = f
}

func (m *Service) makeKeyStore() error {

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	if err := os.MkdirAll(m.keystoreDir, 0700); err != nil {
		return err
	}

	m.keyStore = keystore.NewKeyStore(m.keystoreDir, scryptN, scryptP)

	m.am = accounts.NewManager(m.keyStore)

	return nil
}

func (m *Service) unlockAccounts() error {

	if len(m.keyStore.Accounts()) == 0 {
		return nil
	}

	pwd, err := m.readPwd()
	if err != nil {
		m.logger.WithError(err).Error("Reading PwdFile")
		return err
	}

	for _, ac := range m.keyStore.Accounts() {
		if err := m.keyStore.Unlock(ac, string(pwd)); err != nil {
			return err
		}
		m.logger.WithField("address", ac.Address.Hex()).Debug("Unlocked account")
	}
	return nil
}

func (m *Service) createGenesisAccounts() error {
	if _, err := os.Stat(m.genesisFile); os.IsNotExist(err) {
		return nil
	}

	contents, err := ioutil.ReadFile(m.genesisFile)
	if err != nil {
		return err
	}

	var genesis struct {
		Alloc common.AccountMap
	}

	if err := json.Unmarshal(contents, &genesis); err != nil {
		return err
	}

	if err := m.state.CreateAccounts(genesis.Alloc); err != nil {
		return err
	}
	return nil
}

func (m *Service) serveAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/account/{address}", m.makeHandler(accountHandler)).Methods("GET")
	r.HandleFunc("/accounts", m.makeHandler(accountsHandler)).Methods("GET")
	r.HandleFunc("/block/{hash}", m.makeHandler(blockByHashHandler)).Methods("GET")
	r.HandleFunc("/blockById/{id}", m.makeHandler(blockByIdHandler)).Methods("GET")
	//r.HandleFunc("/blockIndex", m.makeHandler(blockIndexHandler)).Methods("GET")
	r.HandleFunc("/call", m.makeHandler(callHandler)).Methods("POST")
	r.HandleFunc("/tx", m.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/transactions", m.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/rawtx", m.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/sendRawTransaction", m.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/tx/{tx_hash}", m.makeHandler(txReceiptHandler)).Methods("GET")
	r.HandleFunc("/transaction/{tx_hash}", m.makeHandler(transactionReceiptHandler)).Methods("GET")
	r.HandleFunc("/info", m.makeHandler(infoHandler)).Methods("GET")
	r.HandleFunc("/html/info", m.makeHandler(htmlInfoHandler)).Methods("GET")
	http.Handle("/", &CORSServer{r})
	if err := http.ListenAndServe(m.apiAddr, nil); err != nil {
		panic(err)
	}
}

type CORSServer struct {
	r *mux.Router
}

func (s *CORSServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func (m *Service) makeHandler(fn func(http.ResponseWriter, *http.Request, *Service)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		fn(w, r, m)
		m.Unlock()
	}
}

func (m *Service) checkErr(err error) {
	if err != nil {
		m.logger.WithError(err).Error("ERROR")
		os.Exit(1)
	}
}

func (m *Service) readPwd() (pwd string, err error) {
	text, err := ioutil.ReadFile(m.pwdFile)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines[0], nil
}

func (s *Service) AccountManager() *accounts.Manager {
	return s.am
}

func (m *Service) ChainConfig() *params.ChainConfig {
	return m.chainConfig
}
