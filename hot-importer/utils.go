package main

import (
	"flag"
	"os"
	"path/filepath"
)

// Config has all the options you defined at the command line.
type Config struct {
	BootnodesPath    string
	NodeDatabasePath string
	Verbosity        int
	Vmodule          string
}

// ParseFlags gets those command line options and set them into a nice
// Config struct.
func ParseFlags() *Config {
	c := &Config{}

	flag.StringVar(&c.BootnodesPath, "bootnodes", "", "Location of bootnodes file")
	flag.StringVar(&c.NodeDatabasePath, "nodes-database", "", "Location of the node database")
	flag.IntVar(&c.Verbosity, "verbosity", 3, "overall verbosity of the devp2p modules")
	flag.StringVar(&c.Vmodule, "vmodule", "devp2p=5", "verbosity per module, i.e devp2p=5,p2p=3")
	flag.Parse()

	if c.NodeDatabasePath == "" {
		homeDir := os.Getenv("HOME")
		c.NodeDatabasePath = filepath.Join(homeDir, ".mustekala", "devp2p", "nodes")
	}

	return c
}
