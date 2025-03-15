package universities

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/base/data"
	"CertiBlock/application/shared/utils"

	"errors"
)

func SaveCertificateFile(context *base.ApplicationContext, fileUpload *data.CertificateFileUpload) (*string, error) {
	universityPrivateKey := fileUpload.UniversityPrivateKey
	/*universityPublicKey*/ _, err := utils.ComputePublicKeyString(universityPrivateKey)
	if err != nil {
		return nil, errors.New("error computing public key")
	}

	_, err = context.DB.Exec(
		"INSERT INTO certificates (uuid, university_encrypted_ks_1, university_encrypted_ks_2, student_encrypted_ks_1, student_encrypted_ks_2, ks_encrypted_file) VALUES (?, ?, ?, ?, ?, ?)",
		fileUpload.CertUUID,
		fileUpload.UniversityEncryptedKS1,
		fileUpload.UniversityEncryptedKS2,
		fileUpload.StudentEncryptedKS1,
		fileUpload.StudentEncryptedKS2,
		fileUpload.KSEncryptedFile,
	)
	if err != nil {
		return nil, errors.New("error inserting cert file")
	}

	return nil, nil
}
