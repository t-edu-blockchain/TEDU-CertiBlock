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
	"context"
	"fmt"

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
func (a *App) IssueCertificate(certHash, universitySignature, studentSignature, dateOfIssuing, certUUID, universityPK, studentPK string) string {
	result, err := IssueCertificate(a.contract, certHash, universitySignature, studentSignature, dateOfIssuing, certUUID, universityPK, studentPK)
	if err != nil {
		return fmt.Sprintf("Error: Failed to issue certificate: %v", err)
	}
	return fmt.Sprintf("Certificate issued successfully:\n%s", result)
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