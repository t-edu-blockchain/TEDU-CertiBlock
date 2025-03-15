package data

type StudentInput struct {
	CountryID   int    `json:"countryID"`
	NIN         string `json:"NIN"`
	FullName    string `json:"fullName"`
	DateOfBirth string `json:"dateOfBirth"`
	Password    string `json:"password"`
}

type StudentOutput struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	FullName   string `json:"fullName"`
}

type StudentAuth struct {
	PrivateKey string `json:"privateKey"`
}
