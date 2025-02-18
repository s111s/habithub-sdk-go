package main

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
)

const FundAmount = 100_000_000
const TransferAmount = 1_000

func example(networkConfig aptos.NetworkConfig) {

	client, err := aptos.NewClient(networkConfig)
	if err != nil {
		panic("Failed to create client: " + err.Error())
	}

	alice, err := aptos.NewEd25519Account()
	if err != nil {
		panic("Failed to create alice" + err.Error())
	}

	bob, err := aptos.NewEd25519Account()
	if err != nil {
		panic("Failed to create alice " + err.Error())
	}

	fmt.Printf("\n=== Addresses ===\n")
	fmt.Printf("Alice: %s\n", alice.Address.String())
	fmt.Printf("Bob: %s\n", bob.Address.String())

	// Fund the sender with the faucet to create it on-chain
	err = client.Fund(alice.Address, FundAmount)
	if err != nil {
		panic("Failed to fund alice: " + err.Error())
	}

	aliceBalance, err := client.AccountAPTBalance(alice.Address)
	if err != nil {
		panic("Failed to retrieve alice balance: " + err.Error())
	}

	bobBalance, err := client.AccountAPTBalance(bob.Address)
	if err != nil {
		panic("Failed to retrieve bob balance: " + err.Error())
	}

	fmt.Printf("\n=== Initial Balances ===\n")
	fmt.Printf("Alice: %d\n", aliceBalance)
	fmt.Printf("Bob: %d\n", bobBalance)

	// 1. Build Transaction
	accountBytes, err := bcs.Serialize(&bob.Address)
	if err != nil {
		panic("Failed to serialize bob's address: " + err.Error())
	}

	amountBytes, err := bcs.SerializeU64(TransferAmount)
	if err != nil {
		panic("Failed to serialize transfer amount: " + err.Error())
	}

	rawTxn, err := client.BuildTransaction(alice.AccountAddress(), aptos.TransactionPayload{
		Payload: &aptos.EntryFunction{
			Module: aptos.ModuleId{
				Address: aptos.AccountOne,
				Name:    "aptos_account",
			},
			Function: "transfer",
			ArgTypes: []aptos.TypeTag{},
			Args: [][]byte{
				accountBytes,
				amountBytes,
			},
		},
	})

	if err != nil {
		panic("Failed to build transaction: " + err.Error())
	}

	// 2. Simulate Transaction
	simulationResult, err := client.SimulateTransaction(rawTxn, alice)
	if err != nil {
		panic("Failed to simulate transaction: " + err.Error())
	}

	fmt.Printf("\n=== Simulation ===\n")
	fmt.Printf("Gas unit price: %d\n", simulationResult[0].GasUnitPrice)
	fmt.Printf("Gas used: %d\n", simulationResult[0].GasUsed)
	fmt.Printf("Total gas fee: %d\n", simulationResult[0].GasUsed*simulationResult[0].GasUnitPrice)
	fmt.Printf("Status: %s\n", simulationResult[0].VmStatus)

}

func main() {
	example(aptos.DevnetConfig)
}
