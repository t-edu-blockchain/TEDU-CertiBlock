// package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/hyperledger/fabric-gateway/pkg/client"
// 	"google.golang.org/grpc"
// )
// // App struct
// type App struct {
// 	ctx context.Context
// 	gw      *client.Gateway
// 	conn    *grpc.ClientConn
// 	contract *client.Contract
// }

// // NewApp creates a new App application struct
// func NewApp() *App {
// 	return &App{}
// }

// // startup is called when the app starts. The context is saved
// // so we can call the runtime methods
// func (a *App) startup(ctx context.Context) {
// 	a.ctx = ctx
// }

// // Greet returns a greeting for the given name
// func (a *App) Greet(name string) string {
// 	return fmt.Sprintf("Hello %s, It's show time!", name)
// }

// func (a *App) Connect() {
// 	var err error
// 	a.gw, a.conn, err = NewGateway()
// 	if err != nil {
// 		panic(fmt.Errorf("failed to initialize gateway: %w", err))
// 	}
// 	network := a.gw.GetNetwork("mychannel")
// 	a.contract = network.GetContract("certicontract")
// }

// func (a *App) InitLedger(){
// 		// Khởi tạo ledger
// 	result, err := InitLedger(a.contract)
// 	if err != nil {
// 		fmt.Printf("Failed to initialize ledger: %v\n", err)
// 		return
// 	}
// 	fmt.Println("Ledger initialized:", result)
// }

// func (a *App) GetAll() string{
// 	result := QueryAllCertificates(a.ctx, a.contract)
// 	return result
// }

package main

import (
	"CertiBlock/application/shared/utils"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
)

// App struct
type App struct {
	ctx      context.Context
	gw       *client.Gateway
	conn     *grpc.ClientConn
	contract *client.Contract
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Shutdown closes the Gateway and gRPC connections
func (a *App) Shutdown() string {
	if a.gw != nil {
		a.gw.Close()
	}
	if a.conn != nil {
		a.conn.Close()
	}
	return "Disconnected from Fabric Gateway"
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Connect establishes a connection to the Fabric Gateway
func (a *App) Connect() string {
	var err error
	a.gw, a.conn, err = NewGateway()
	if err != nil {
		return fmt.Sprintf("Error: Failed to initialize gateway: %v", err)
	}
	network := a.gw.GetNetwork("mychannel")
	a.contract = network.GetContract("certicontract")
	return "Connected to Fabric Gateway and chaincode 'certicontract' on channel 'mychannel'"
}

// InitLedger initializes the ledger
func (a *App) InitLedger() string {
	result, err := InitLedger(a.contract)
	if err != nil {
		return fmt.Sprintf("Error: Failed to initialize ledger: %v", err)
	}
	return fmt.Sprintf("Ledger initialized:\n%s", result)
}

// IssueCertificate issues a new certificate
func (a *App) IssueCertificate(universitySignature, studentSignature, dateOfIssuing, universityPrivateKeyString, studentPublicKeyString, backendServerUrl string) string {
	ctx := a.ctx
	options := runtime.OpenDialogOptions{
		Title: "Select a file",
	}

	filePath, err := runtime.OpenFileDialog(ctx, options)
	if err != nil {
		return fmt.Sprintf("Error selecting file: %s", err)
	}
	if filePath == "" {
		return "no file selected!"
	}

	// Read into byte slice
	fileRawBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error: Failed to read file: %v", err)
	}

	// Base64 encode the content
	base64File := base64.StdEncoding.EncodeToString(fileRawBytes)

	////////////////////////////////
	//////// END FILE STUFF ////////
	////////////////////////////////

	certUUID := uuid.NewString()

	universityPublicKeyString, err := utils.ComputePublicKeyString(universityPrivateKeyString)
	if err != nil {
		return fmt.Sprintf("Error: Failed to compute university public key: %v", err)
	}
	fmt.Println("EEEEEEE2")

	certHash := utils.HashSHA512(base64File + dateOfIssuing + certUUID + studentPublicKeyString + universityPublicKeyString)

	if a.contract == nil {
		return "Please connect and init ledger first!"
	}

	ks := utils.GenerateSecureRandomString(64)
	studentEncryptedKS1, studentEncryptedKS2, err := utils.ElGamalEncryptString(studentPublicKeyString, ks)
	if err != nil {
		return fmt.Sprintf("Error: Failed to encrypt KS using student pubkey: %v", err)
	}
	fmt.Println("EEEEEEE4")

	universityEncryptedKS1, universityEncryptedKS2, err := utils.ElGamalEncryptString(universityPublicKeyString, ks)
	if err != nil {
		return "Error: Failed to encrypt KS using university pubkey"
	}

	fmt.Println("EEEEEEE5")
	ksEncryptedFile, err := utils.VigenereEncryptString(ks, base64File)
	if err != nil {
		return "Error: Failed to encrypt file using KS"
	}

	encryptedFileSize := len(ksEncryptedFile)
	fmt.Printf("Encrypted file size: %d\n", encryptedFileSize)
	if encryptedFileSize > 66386403 {
		return "Error: File too large. Maximum size is about 26MB. (This is a known issue due to MySQL. Thank you for your understanding!)"
	}

	fmt.Println("EEEEEEE6")
	// Save the file to the backend
	payload := gin.H{
		"certUUID":               certUUID,
		"ksEncryptedFile":        ksEncryptedFile,
		"studentEncryptedKS1":    studentEncryptedKS1,
		"studentEncryptedKS2":    studentEncryptedKS2,
		"universityEncryptedKS1": universityEncryptedKS1,
		"universityEncryptedKS2": universityEncryptedKS2,
		"universityPrivateKey":   universityPrivateKeyString,
	}

	jsonPayload, err := json.Marshal(payload)
	fmt.Println("EEEEEEE7")
	if err != nil {
		return fmt.Sprintf("Error: Failed to marshal payload: %v", err)
	}

	fmt.Println("HERE")
	res, err := http.Post(backendServerUrl+"/api/universities/certificate-file", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Sprintf("Error: Failed to save file to backend: %v", err)
	}
	defer res.Body.Close()

	fmt.Println("HERE2")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Sprintf("Error: Failed to read response body: %v", err)
	}

	fmt.Println("HERE3")

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return fmt.Sprintf("Error: Failed to save file to backend: %s", string(body))
	}

	fmt.Println("HERE4")

	result, err := IssueCertificate(a.contract, certHash, universitySignature, studentSignature, dateOfIssuing, certUUID, universityPublicKeyString, studentPublicKeyString)
	if err != nil {
		return fmt.Sprintf("Error: Failed to issue certificate: %v", err)
	}
	fmt.Println("EEEEEEE3")

	return fmt.Sprintf("Certificate issued successfully:\n%s\n\nUniversity Public Key:%s\n", result, universityPublicKeyString)
}

// RegisterUniversity registers a new university
func (a *App) RegisterUniversity(name, publicKey, location, description string) string {
	result, err := RegisterUniversity(a.contract, name, publicKey, location, description)
	if err != nil {
		return fmt.Sprintf("Error: Failed to register university: %v", err)
	}
	return fmt.Sprintf("University registered successfully:\n%s", result)
}

// GetAll queries all certificates
func (a *App) GetAll() string {
	result := QueryAllCertificates(a.ctx, a.contract)
	if result == "Looix" {
		return "Error: Failed to query certificates"
	}
	return result
}
