package enrollment_certificates

import (
	"certiblock/base"
	"certiblock/base/data"
	"certiblock/services/universities"
)

func Issue(context *base.ApplicationContext, enrollmentCertificateInput data.BCEnrollmentCertificateInput) (*data.BCEnrollmentCertificateOutput, error) {
	university, err := universities.GetByPrivateKey(context, enrollmentCertificateInput.UniversityPrivateKey)
	if err != nil {
		return nil, err
	}
	result, err := context.DB.Exec("INSERT INTO bc_enrollment_certificates (university_public_key, student_public_key, hash) VALUES (?, ?, ?)",
		university.PublicKey,
		enrollmentCertificateInput.StudentPublicKey,
		enrollmentCertificateInput.Hash,
	)

	if err != nil {
		return nil, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	output := data.BCEnrollmentCertificateOutputResponse(data.BCEnrollmentCertificate{
		ID:                     int(lastInsertedID),
		StudentPublicKey:       enrollmentCertificateInput.StudentPublicKey,
		GetUniversityPublicKey: university.PublicKey,
		Hash:                   enrollmentCertificateInput.Hash,
	})
	return &output, nil
}
