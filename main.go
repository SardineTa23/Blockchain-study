package main

import (
	"blockchain-study/wallet"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	// myBlockChainAddress := "my_blockchain_address"
	// blockchain := NewBlockChain(myBlockChainAddress)
	// blockchain.Print()

	// blockchain.AddTransaction("A", "B", 1.0)
	// blockchain.Mining()
	// blockchain.Print()

	// blockchain.AddTransaction("C", "D", 2.0)
	// blockchain.AddTransaction("X", "Y", 3.0)
	// blockchain.Mining()
	// blockchain.Print()

	// fmt.Printf("my %.1f\n", blockchain.CalculateTotalAmount(myBlockChainAddress))
	// fmt.Printf("C %.1f\n", blockchain.CalculateTotalAmount("C"))
	// fmt.Printf("D %.1f\n", blockchain.CalculateTotalAmount("D"))

	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
	fmt.Println(w.BlockchainAddress())

	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "B", 1.0)
	fmt.Printf("signature %s \n", t.GenerateSignature())
}
