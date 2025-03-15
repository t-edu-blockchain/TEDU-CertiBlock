package data

type CertificateOutput struct {
	CertHash            string `json:"certHash"`
	CertUUID            string `json:"certUUID"`
	StudentPublicKey    string `json:"studentPublicKey"`
	UniversityPublicKey string `json:"universityPublicKey"`
	DateOfIssuing       string `json:"dateOfIssuing"`
}

type CertificateOutputFull struct {
	CertificateOutput
	Base64File string `json:"base64File"`
}

// Certificate as stored in the blockchain
type BlockchainCertificateOutput struct {
	CertHash            string `json:"certHash"`
	UniversitySignature string `json:"universitySignature"`
	StudentSignature    string `json:"studentSignature"`
	DateOfIssuing       string `json:"dateOfIssuing"`
	CertUUID            string `json:"certUUID"`
	UniversityPublicKey string `json:"universityPK"`
	StudentPublicKey    string `json:"studentPK"`
}

// Certificate file upload from university
type CertificateFileUpload struct {
	CertUUID               string `json:"certUUID"`
	UniversityPrivateKey   string `json:"universityPrivateKey"`
	UniversityEncryptedKS1 string `json:"universityEncryptedKS1"`
	UniversityEncryptedKS2 string `json:"universityEncryptedKS2"`
	StudentEncryptedKS1    string `json:"studentEncryptedKS1"`
	StudentEncryptedKS2    string `json:"studentEncryptedKS2"`
	KSEncryptedFile        string `json:"ksEncryptedFile"`
}
