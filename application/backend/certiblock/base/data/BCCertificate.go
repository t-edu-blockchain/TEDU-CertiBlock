package data

type BCCertificateInput struct {
	File                 string `json:"file"`
	StudentPublicKey     string `json:"student_public_key"`
	UniversityPrivateKey string `json:"university_private_key"`
}

type BCCertificateInputValidated struct {
	V_h                  string `json:"hash"`
	V_C_X                string `json:"C_X"`
	V_C_U                string `json:"C_U"`
	V_C_S                string `json:"C_S"`
	StudentPublicKey     string `json:"student_public_key"`
	UniversityPrivateKey string `json:"university_private_key"`
	V_S0                 string `json:"S0"`
	V_S                  string `json:"S"`
}

type BCCertificate struct {
	ID int `json:"id"`
	BCCertificateInputValidated
}

type BCCertificateOutput struct {
	BCCertificate
}

func BCCertificateOutputResponse(bcCertificate any) BCCertificateOutput {
	switch c := bcCertificate.(type) {
	case BCCertificate:
		return BCCertificateOutput{
			BCCertificate: c,
		}
	case *BCCertificate:
		return BCCertificateOutputResponse(*c)
	case BCCertificateOutput:
		return c
	case *BCCertificateOutput:
		return *c
	default:
		// Return an empty BCCertificateOutput if conversion is not possible
		return BCCertificateOutput{}
	}
}
