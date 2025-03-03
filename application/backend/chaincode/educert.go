package main

import (
    "CertiBlock/chaincode/educert"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
    chaincode, err := contractapi.NewChaincode(&educert.SmartContract{})
    if err != nil {
        panic("Error creating educert chaincode: " + err.Error())
    }

    if err := chaincode.Start(); err != nil {
        panic("Error starting educert chaincode: " + err.Error())
    }
}