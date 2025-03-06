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

// InitLedger khởi tạo ledger với schema ban đầu
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
		"issueCertificate",
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
		"registerUniversity",
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

// QueryCertificateByUUID truy vấn chứng chỉ theo UUID
func QueryCertificateByUUID(contract *client.Contract, certUUID string) (string, error) {
	result, err := contract.EvaluateTransaction("queryCertificateByUUID", certUUID)
	if err != nil {
		return "", fmt.Errorf("failed to query certificate with UUID %s: %w", certUUID, err)
	}
	return formatJSON(result), nil
}

// GetAllCertificatesByStudent lấy tất cả chứng chỉ của một sinh viên
func GetAllCertificatesByStudent(contract *client.Contract, studentPK string) (string, error) {
	result, err := contract.EvaluateTransaction("getAllCertificateByStudent", studentPK)
	if err != nil {
		return "", fmt.Errorf("failed to get certificates for student %s: %w", studentPK, err)
	}
	return formatJSON(result), nil
}

// GetAllCertificatesByUniversity lấy tất cả chứng chỉ của một trường đại học
func GetAllCertificatesByUniversity(contract *client.Contract, universityPK string) (string, error) {
	result, err := contract.EvaluateTransaction("getAllCertificateByUniversity", universityPK)
	if err != nil {
		return "", fmt.Errorf("failed to get certificates for university %s: %w", universityPK, err)
	}
	return formatJSON(result), nil
}

// QueryAllCertificates lấy tất cả chứng chỉ
func QueryAllCertificates(contract *client.Contract) (string, error) {
	result, err := contract.EvaluateTransaction("queryAll")
	if err != nil {
		return "", fmt.Errorf("failed to query all certificates: %w", err)
	}
	return formatJSON(result), nil
}

// formatJSON định dạng JSON cho dễ đọc
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return string(data) // Trả về nguyên gốc nếu lỗi
	}
	return prettyJSON.String()
}