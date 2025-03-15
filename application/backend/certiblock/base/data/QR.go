package data

type QRInput struct {
	StudentPrivateKey string `json:"studentPrivateKey"`
	CertUUID          string `json:"certUUID"`
}

type QROutput struct {
	Content string `json:"content"`
}
