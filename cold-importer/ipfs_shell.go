package main

import (
	api "github.com/ipfs/go-ipfs-api"
)

func initStart() *api.Shell {
	ipfsapi := "locahost:5001"
	ipfs := api.NewShell(ipfsapi)

	for i := 0; !ipfs.IsUp(); i++ {
		// log.Println("ipfs daemon not running, waiting a second...")
		// un time sleep acá, 100ms
	}

	return ipfs
}

// Funcion para saber Si tengo ya ese cid en mi storage

// Funcion para insertar dicho cid en el storage
