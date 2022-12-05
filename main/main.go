package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

var (
	cfg    *Config
	client *ethclient.Client
)

func main() {
	cfg = LoadConfig()
	var err error
	client, err = ethclient.Dial(cfg.RpcUrl)

	if err != nil {
		log.Panicf("dial %s failed, err=%v\n", cfg.RpcUrl, err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Panicf("get chainid failed, err=%v\n", err)
	}
	log.Println("chainID ", chainID)

	auth, _ := bind.NewKeyedTransactorWithChainID(cfg.secret, chainID)
	// if the account will send transactions concurrently,
	// you need to manage the nonce by yourself
	// auth.Nonce = manualSetNonce

	auth.GasTipCap = big.NewInt(1 * params.GWei)
	// deploy
	_, tx, contract, err := DeployEIP20(auth, client, big.NewInt(1e18), "dtoken", 8, "dt")
	if err != nil {
		log.Panicf("DeployEIP20 failed, err=%v\n", err)
	}

	var receipt *types.Receipt
	if cfg.isHttp {
		receipt, err = WaitReceipt(client, tx.Hash(), 3*time.Second, 5, 3)
	} else {
		receipt, err = WaitReceiptOnNewHead(client, tx.Hash(), 5, 3)
	}
	if err != nil {
		log.Panicf("WaitReceiptOnNewHead failed, err=%v\n", err)
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		log.Printf("contract is deployed at %s", receipt.ContractAddress.String())
	} else {
		strres, _ := json.Marshal(receipt)
		log.Println(string(strres))
		log.Panic("deploy contract failed")
	}

	// transact
	toAddr := common.Address{0x11}
	tx, err = contract.Transfer(auth, toAddr, big.NewInt(12345))
	if err != nil {
		log.Panicf("Transfer failed, err=%v\n", err)
	}
	log.Printf("Send erc20 transfer tx, hash=%s\n", tx.Hash())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	if cfg.isHttp {
		go FilterTransferEvent(wg, contract, client, 3)
	} else {
		go WatchTransferEvent(wg, contract)
	}
	wg.Wait()
	log.Println("Done")
}

func WaitReceipt(conn *ethclient.Client, txHash common.Hash, tryInterval time.Duration, tryTimes int, confirms int) (*types.Receipt, error) {
	ticker := time.NewTicker(tryInterval)
	defer ticker.Stop()
	i := 0
	for {
		select {
		case <-ticker.C:
			i++
			r, err := conn.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				if i >= tryTimes {
					return nil, err
				}
				log.Printf("can't get the receipt, err: %v\nwill try again.\n", err)
			} else {
				head, err := conn.BlockNumber(context.Background())
				if err != nil {
					return nil, err
				}
				confirm := head - r.BlockNumber.Uint64()
				log.Printf("got receipt, txHash=%x, confirms=%d \n", txHash, confirm)
				if confirm >= uint64(confirms) {
					return r, nil
				}
			}
		}
	}
}

// If the rpc-server supports subscription,
// then we can check the receipt only when there's a new block.
func WaitReceiptOnNewHead(conn *ethclient.Client, txHash common.Hash, waitBlocks int, confirms int) (*types.Receipt, error) {
	headers := make(chan *types.Header)
	sub, err := conn.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal("Error to SubscribeNewHead:", err)
	}
	defer sub.Unsubscribe()
	for i := 0; i < waitBlocks+confirms; i++ {
		select {
		case err := <-sub.Err():
			return nil, err
		case header := <-headers:
			r, err := conn.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				log.Printf("can't get the receipt, err: %v\nwill try again.\n", err)
			} else {
				confirm := header.Number.Uint64() - r.BlockNumber.Uint64()
				log.Printf("got receipt, txHash=%x, confirms=%d \n", txHash, confirm)
				if confirm >= uint64(confirms) {
					return r, nil
				}
			}
		}
	}
	return nil, errors.New("not found")
}

func WatchTransferEvent(wg *sync.WaitGroup, contract *EIP20) {
	defer wg.Done()
	sink := make(chan *EIP20Transfer)
	sub, err := contract.WatchTransfer(new(bind.WatchOpts), sink, []common.Address{}, []common.Address{})
	if err != nil {
		log.Printf("WatchTransfer failed: %v\n", err)
		return
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			if err != nil {
				log.Println("subscribe return: ", err)
			}
			return
		case transfer := <-sink:
			log.Printf("Got ERC20 transfer event at Block %d,\ntxHash: %s, from: %s, to: %s, value: %d\n",
				transfer.Raw.BlockNumber, transfer.Raw.TxHash, transfer.From, transfer.To, transfer.Value)
			log.Println("Stop watching")
			sub.Unsubscribe()
		}
	}
}

func FilterTransferEvent(wg *sync.WaitGroup, contract *EIP20, client *ethclient.Client, confirms int) {
	defer wg.Done()

	var start uint64
	for {
		num, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Printf("Get BlockNumber failed: %v\n", err)
			return
		}
		if num > uint64(confirms) {
			start = num - uint64(confirms)
			break
		}
		time.Sleep(3 * time.Second)
	}
	for {
		isOk := false
		end := start + 1
		opts := bind.FilterOpts{Start: start, End: &end}
		it, err := contract.FilterTransfer(&opts, []common.Address{}, []common.Address{})
		if err != nil {
			log.Printf("FilterTransfer failed: %v\n", err)
			return
		}

		for it.Next() {
			transfer := it.Event
			log.Printf("Got ERC20 transfer event at Block %d,\ntxHash: %s, from: %s, to: %s, value: %d\n",
				transfer.Raw.BlockNumber, transfer.Raw.TxHash, transfer.From, transfer.To, transfer.Value)
			log.Println("Close filter")
			it.Close()
			isOk = true
		}
		if isOk {
			return
		}
		start++ // next block
		time.Sleep(3 * time.Second)
	}

}
