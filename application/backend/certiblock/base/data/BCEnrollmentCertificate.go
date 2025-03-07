package data

type BCEnrollmentCertificateInput struct {
	StudentPublicKey     string `json:"student_public_key"`
	UniversityPrivateKey string `json:"university_private_key"`
	Hash                 string `json:"hash"`
}

type BCEnrollmentCertificate struct {
	ID                     int    `json:"id"`
	StudentPublicKey       string `json:"student_public_key"`
	GetUniversityPublicKey string `json:"university_public_key"`
	Hash                   string `json:"hash"`
}

type BCEnrollmentCertificateOutput struct {
	BCEnrollmentCertificate
}

func BCEnrollmentCertificateOutputResponse(bcEnrollmentCertificate any) BCEnrollmentCertificateOutput {
	switch c := bcEnrollmentCertificate.(type) {
	case BCEnrollmentCertificate:
		return BCEnrollmentCertificateOutput{
			BCEnrollmentCertificate: c,
		}
	case *BCEnrollmentCertificate:
		return BCEnrollmentCertificateOutput{
			BCEnrollmentCertificate: *c,
		}
	case BCEnrollmentCertificateOutput:
		return c
	case *BCEnrollmentCertificateOutput:
		return *c
	default:
		// Return an empty BCEnrollmentCertificateOutput if conversion is not possible
		return BCEnrollmentCertificateOutput{}
	}
}
