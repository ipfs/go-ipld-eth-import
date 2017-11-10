package devp2p

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	colorable "github.com/mattn/go-colorable"
)

// Manager is the object that organizes and enables the connection
// of a node to the devp2p network, it can start and stop a server
// and make available data and metrics through an API.
type Manager struct {
	server *p2p.Server
}

// Config is the configuration object for DevP2P
type Config struct {
	// bootnodes file
	BootnodesPath string

	// bootnodes slice
	bootnodes []*discover.Node

	// node database path. Must be appointed outside this package
	NodeDatabasePath string

	// we can find the client's private key here
	PrivateKeyFilePath string

	// glogger verbosity level (5 is the highest)
	Verbosity int

	// glogger verbosity per module (ex: devp2p=5,p2p=5)
	Vmodule string
}

// NewDevP2P returns a DevP2P Manager object
//
// * defines logger.
//
// * defines bootnodes.
//
// * defines node database.
//
// * sets up the _peerstore_.
//
// * defines and configures _server_, passing the _protocol-handler_.
// * _protocol-handler_ needs the _peer-status-msg_ methods to perform
// * an eth handshake, adding and removing peers from the _peerstore_.
// * also, _protocol-handler_ will put some received requests into channels,
// * we may want to answer them.
//
// * sets up _api_, this will talk to the _peer-send_ methods,
// * which in turn picks the best peers and ask the question.
//
// * sets up the _metrics_.
//
func NewManager(config Config) *Manager {
	var err error

	manager := &Manager{}

	setupLogger(config)

	config.bootnodes, err = parseBootnodesFile(config.BootnodesPath)
	if err != nil {
		log.Error("NewManager", "error", fmt.Sprintf("processBootnodesFile error: %v", err))
		os.Exit(1) // zero tolerance
	}

	if config.NodeDatabasePath == "" {
		log.Error("NewManager", "error", "node database path must be appointed outside this package")
		os.Exit(1)
	}

	// TODO
	// Setup the peerstore

	manager.server = newServer(config)

	// TODO
	// setup API

	// TODO
	// setup metrics

	return manager
}

// Start should be run as a goroutine
func (m *Manager) Start() {
	if err := m.server.Start(); err != nil {
		log.Error("error starting devp2p server", "error", err)
		os.Exit(1)
	}
}

// Stop terminates the server
func (m *Manager) Stop() {
	m.server.Stop()
}

// setupLogger configures glogger with the required verbosity.
func setupLogger(config Config) {
	output := colorable.NewColorableStderr()
	glogger := log.NewGlogHandler(log.StreamHandler(output, log.TerminalFormat(true)))

	glogger.Verbosity(log.Lvl(config.Verbosity))
	glogger.Vmodule(config.Vmodule)

	log.Root().SetHandler(glogger)
}

// parseBootnodesFile parses the bootnodes file to be included in the
// devp2p server.
func parseBootnodesFile(filePath string) ([]*discover.Node, error) {
	nodes := []*discover.Node{}

	if filePath == "" {
		return nil, fmt.Errorf("A bootnodes file must be defined!")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nodeUrl := scanner.Text()
		node, err := discover.ParseNode(nodeUrl)
		if err != nil {
			log.Error("add bootstrap node error", "node-url", nodeUrl, "error", err)
		}
		nodes = append(nodes, node)
		log.Debug("added bootstrap node", "node-url", nodeUrl)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}
