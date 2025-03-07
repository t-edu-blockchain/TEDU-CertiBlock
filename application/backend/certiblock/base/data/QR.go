package data

type QRInput struct {
	V_rk string `json:"rk"`
	V_h  string `json:"h"`
	// StudentPublicKey string `json:"student_public_key"`
}

type QR struct {
	ID int `json:"id"`
	QRInput
}

type QROutput struct {
	Token string `json:"token"`
}
