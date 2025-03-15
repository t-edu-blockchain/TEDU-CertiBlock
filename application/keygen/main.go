package main

import (
	"CertiBlock/application/shared/utils"
	"fmt"
)

func main() {
	privateKeyString, err := utils.GeneratePrivateKeyString()
	if err != nil {
		fmt.Printf("Error generating private key: %v\n", err)
		return
	}
	publicKeyString, err := utils.ComputePublicKeyString(privateKeyString)
	if err != nil {
		fmt.Printf("Error computing public key: %v\n", err)
		return
	}
	fmt.Printf("Private key: %s\n", privateKeyString)
	fmt.Printf("Public key: %s\n", publicKeyString)
}
