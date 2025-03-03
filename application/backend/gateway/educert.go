// gateway/educert.go
package educert

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "UETMSP"
	cryptoPath   = "../../../network/organizations/peerOrganizations/uet"
	certPath     = cryptoPath + "/users/User1@uet/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@uet/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.uet/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.uet"
)

// NewGateway khởi tạo kết nối tới Fabric Gateway
func NewGateway() (*client.Gateway, *grpc.ClientConn, error) {
	clientConnection := newGrpcConnection()
	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, nil, err
	}
	return gw, clientConnection, nil
}

func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

// Các hàm tương tác với chaincode
func InitLedger(contract *client.Contract) error {
	_, err := contract.SubmitTransaction("InitLedger")
	return err
}

func CreateCertificate(contract *client.Contract, id, studentName, university, issueDate, degreeType string) error {
	_, err := contract.SubmitTransaction("CreateCertificate", id, studentName, university, issueDate, degreeType)
	return err
}

func ReadCertificate(contract *client.Contract, id string) (string, error) {
	result, err := contract.EvaluateTransaction("ReadCertificate", id)
	if err != nil {
		return "", err
	}
	return formatJSON(result), nil
}

func UpdateCertificate(contract *client.Contract, id, studentName, university, issueDate, degreeType string) error {
	_, err := contract.SubmitTransaction("UpdateCertificate", id, studentName, university, issueDate, degreeType)
	return err
}

func GetAllCertificates(contract *client.Contract) (string, error) {
	result, err := contract.EvaluateTransaction("GetAllCertificates")
	if err != nil {
		return "", err
	}
	return formatJSON(result), nil
}

func DeleteCertificate(contract *client.Contract, id string) error {
	_, err := contract.SubmitTransaction("DeleteCertificate", id)
	return err
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return string(data) // Trả về nguyên gốc nếu lỗi
	}
	return prettyJSON.String()
}