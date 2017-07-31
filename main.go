package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	rlp "github.com/ethereum/go-ethereum/rlp"

	api "github.com/ipfs/go-ipfs-api"
)

func getEthClient(ethapi string) *ethclient.Client {
	for i := 0; true; i++ {
		ethcli, err := ethclient.Dial(ethapi)
		if err == nil {
			return ethcli
		}

		if i > 30 {
			log.Fatal("timed out waiting for ethereum daemon to come online")
		}

		if !(strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "bad status")) {
			log.Fatal("unexpected error waiting on geth: ", err)
		}

		log.Println("ethereum daemon not running, trying again in a second...")
		time.Sleep(time.Second)
	}
	panic("should never reach this")
}

func handleNewBlock(cli *ethclient.Client, ipfs *api.Shell, h common.Hash) {
	blk, err := cli.BlockByHash(context.Background(), h)
	if err != nil {
		log.Fatal("error getting block:", err)
	}

	blkdata, err := rlp.EncodeToBytes(blk)
	if err != nil {
		log.Fatal("error rlp encoding block:", err)
	}

	cid, err := ipfs.DagPut(blkdata, "raw", "eth")
	if err != nil {
		log.Fatal("error from dag put:", err)
	}

	fmt.Println("new block: ", cid)
}

func main() {
	ipfsapi := flag.String("ipfsapi", "localhost:5001", "ipfs api host and port")
	ethapi := flag.String("ethapi", "ws://localhost:8546", "ethereum websockets api endpoint")
	flag.Parse()

	cli := getEthClient(*ethapi)

	ipfs := api.NewShell(*ipfsapi)
	for i := 0; !ipfs.IsUp(); i++ {
		log.Println("ipfs daemon not running, waiting a second...")
	}

	if len(flag.Args()) > 0 {
		h := common.HexToHash(flag.Arg(0))
		handleNewBlock(cli, ipfs, h)
		return
	}

	nblocks := make(chan *types.Header)
	subs, err := cli.SubscribeNewHead(context.Background(), nblocks)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Subscription started, listening for new blocks")
	for {
		select {
		case h := <-nblocks:
			log.Println("new hash:", h.Hash().Hex())
			handleNewBlock(cli, ipfs, h.Hash())
		case err := <-subs.Err():
			log.Println("error from subscription: ", err)
			log.Println("waiting 10 seconds and retrying")
			time.Sleep(time.Second * 10)
			cli = getEthClient(*ethapi)
			nsubs, err := cli.SubscribeNewHead(context.Background(), nblocks)
			if err != nil {
				log.Fatal("error trying to resubscribe:", err)
			}

			subs = nsubs
		}
	}
}
