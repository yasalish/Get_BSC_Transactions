package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Transaction struct {
	Type         uint8  `json:"Type"`
	Nonce        uint64 `json:"Nonce"`
	Txhash       string `json:"Txhash"`
	Blockno      uint64 `json:"Blockno"`
	To           string `json:"To"`
	From         string `json:"From"`
	CurrentValue string `json:"CurrentValue"`
	Gas          uint64 `json:"Gas"`
	GasPrice     uint64 `json:"GasPrice"`
	Cost         uint64 `json:"Cost"`
	Data         []byte `json:"Data"`
}

func main() {
	http.HandleFunc("/", showTransactions)
	http.ListenAndServe(":8080", nil)
}

func showTransactions(resp http.ResponseWriter, req *http.Request) {

	var ConAddr = "0x1da200f724b6e707cD8B8593f2c270771B7FC769"
	var StartBlock int64 = 11858824
	var EndBlock int64 = 12392042
	var step int64

	var oneTransaction Transaction
	var allTransactions []Transaction

	node, err := ethclient.Dial("https://bsc-dataseed1.defibit.io")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to BSCSCAN is created")
	var num int = 0

	for step = StartBlock; step <= EndBlock; step++ {
		blockNumber := big.NewInt(step)

		block, err := node.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			log.Fatal(err)
		}

		chID, err := node.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		for _, trans := range block.Transactions() {
			if trans.To() != nil {
				if trans.To().Hex() == ConAddr {
					num = num + 1
					fmt.Println(num)

					oneTransaction.Type = trans.Type()
					oneTransaction.Txhash = trans.Hash().Hex() // 0x5d49fcaa394c97ec8a9c3e7bd9e838
					oneTransaction.Blockno = block.Number().Uint64()
					oneTransaction.CurrentValue = trans.Value().String()
					oneTransaction.Gas = trans.Gas()
					oneTransaction.GasPrice = trans.GasPrice().Uint64()
					oneTransaction.Nonce = trans.Nonce()
					oneTransaction.To = trans.To().Hex()
					msg, _ := trans.AsMessage(types.NewEIP155Signer(chID), nil)
					oneTransaction.From = msg.From().Hex()
					oneTransaction.Data = trans.Data()
					oneTransaction.Cost = trans.Cost().Uint64()
					allTransactions = append(allTransactions, oneTransaction)

					fmt.Println(trans.Hash().Hex())
					fmt.Println(block.Number())
					fmt.Println(trans.Value().String())
					fmt.Println(trans.Gas())
					fmt.Println(trans.GasPrice().Uint64())
					fmt.Println(trans.Nonce())
					fmt.Println(trans.Data())
					fmt.Println(trans.To().Hex())

				}
			}

		}
	}

	jsonTXdata, err := json.MarshalIndent(allTransactions, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	jsonTXFile, err := os.Create("./tranactions.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonTXFile.Close()
	jsonTXFile.Write(jsonTXdata)
	jsonTXFile.Close()
	fmt.Println(string(jsonTXdata))

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonTXdata)
	return
}
