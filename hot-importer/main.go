package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

func main() {
	//
	bootnodesPath := flag.String("bootnodes", "", "Location of bootnodes file")
	flag.Parse()

	// Setup the bootstrap nodes
	if *bootnodesPath == "" {
		log.Fatalf("A bootnodes file must be defined!")
	}
	bootnodes, err := getBootnodes(*bootnodesPath)
	if err != nil {
		log.Fatalf("could not read bootstrap nodes file: %v", err)
	}

	//
	homeDir := os.Getenv("HOME")
	nodeDatabase := filepath.Join(homeDir, ".mustekala", "devp2p")

	//
	privKeyString := "af8bf8bb4c634b8716880aa44e82da72b902144940a56e1fa787505aa513ba46"
	privateKey, err := crypto.HexToECDSA(privKeyString)
	if err != nil {
		log.Fatalf("ERROR: ", err)
		os.Exit(1)
	}

	dialer := p2p.TCPDialer{&net.Dialer{Timeout: 60 * time.Second}}

	// Build a new p2pserver
	devp2pConfig := p2p.Config{
		Name:            "Mustekala libp2p devp2p bridge",
		NodeDatabase:    nodeDatabase,
		PrivateKey:      privateKey,
		NoDiscovery:     false,
		BootstrapNodes:  bootnodes,
		ListenAddr:      ":20000",
		Dialer:          dialer,
		MaxPeers:        100000,
		MaxPendingPeers: 100000,
	}
	devp2pServer := &p2p.Server{Config: devp2pConfig}
	log.Println("instance:", devp2pConfig.Name)

	//
	if err := devp2pServer.Start(); err != nil {
		log.Fatalf("%v", err)
	}

	// INFINITE
	select {}
}

func getBootnodes(filePath string) ([]*discover.Node, error) {
	nodes := []*discover.Node{}

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
			log.Fatalf("Bootstrap URL %s: %v\n", nodeUrl, err)
		}
		nodes = append(nodes, node)
		log.Printf("Added Bootstrap Node: %v", nodeUrl)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}
