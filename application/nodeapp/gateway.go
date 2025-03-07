package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
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
	mspID        = "HUSMSP"
	cryptoPath   = "../../network/organizations/peerOrganizations/hus"
	certPath     = cryptoPath + "/users/User1@hus/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@hus/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.hus/tls/ca.crt"
	peerEndpoint = "dns:///localhost:9051"
	gatewayPeer  = "peer0.hus"
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
		return nil, nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}
	return gw, clientConnection, nil
}

// newGrpcConnection tạo kết nối gRPC với peer
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certificate file at %s: %w", tlsCertPath, err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse TLS certificate: %w", err))
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection to %s: %w", peerEndpoint, err))
	}

	return connection
}

// newIdentity tạo định danh X509 cho client
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file from %s: %w", certPath, err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse certificate: %w", err))
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(fmt.Errorf("failed to create X509 identity: %w", err))
	}

	return id
}

// newSign tạo hàm ký giao dịch
func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key from %s: %w", keyPath, err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse private key: %w", err))
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(fmt.Errorf("failed to create signing function: %w", err))
	}

	return sign
}

// readFirstFile đọc file đầu tiên trong thư mục
func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory %s: %w", dirPath, err)
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

func InitLedger(contract *client.Contract) (string, error) {
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		return "", fmt.Errorf("failed to initialize ledger: %w", err)
	}
	return formatJSON(result), nil
}

// IssueCertificate phát hành chứng chỉ mới
func IssueCertificate(contract *client.Contract, certHash, universitySignature, studentSignature, dateOfIssuing, certUUID, universityPK, studentPK string) (string, error) {
	result, err := contract.SubmitTransaction(
		"IssueCertificate",
		certHash,
		universitySignature,
		studentSignature,
		dateOfIssuing,
		certUUID,
		universityPK,
		studentPK,
	)
	if err != nil {
		return "", fmt.Errorf("failed to issue certificate: %w", err)
	}
	return formatJSON(result), nil
}

// RegisterUniversity đăng ký một trường đại học
func RegisterUniversity(contract *client.Contract, name, publicKey, location, description string) (string, error) {
	result, err := contract.SubmitTransaction(
		"RegisterUniversity",
		name,
		publicKey,
		location,
		description,
	)
	if err != nil {
		return "", fmt.Errorf("failed to register university: %w", err)
	}
	return formatJSON(result), nil
}

func QueryAllCertificates(ctx context.Context, contract *client.Contract) string {
    log.Println("Bắt đầu truy vấn QueryAll")
	result, err := contract.EvaluateWithContext(ctx, "QueryAll")
    if err != nil {
        log.Printf("Lỗi truy vấn: %v", err)
        return "Lỗi"
    }
    log.Println("Truy vấn thành công")
    return formatJSON(result)
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return string(data) // Trả về nguyên gốc nếu lỗi
	}
	return prettyJSON.String()
}