package certificates

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/base/data"
	educert "CertiBlock/application/backend/gateway"
	"CertiBlock/application/shared/utils"
	"encoding/json"
)

func GetAllCertificatesByStudent(context *base.ApplicationContext, studentPrivateKeyString string) (*[]data.CertificateOutput, error) {
	studentPublicKeyString, err := utils.ComputePublicKeyString(studentPrivateKeyString)
	if err != nil {
		return nil, err
	}

	result, err := educert.GetAllCertificatesByStudent(context.Contract, studentPublicKeyString)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		list := []data.CertificateOutput{}
		return &list, nil
	}

	var certificates []data.BlockchainCertificateOutput
	err = json.Unmarshal([]byte(result), &certificates)
	if err != nil {
		return nil, err
	}

	var output []data.CertificateOutput
	for _, cert := range certificates {
		output = append(output, data.CertificateOutput{
			CertHash:            cert.CertHash,
			CertUUID:            cert.CertUUID,
			StudentPublicKey:    cert.StudentPublicKey,
			UniversityPublicKey: cert.UniversityPublicKey,
			DateOfIssuing:       cert.DateOfIssuing,
		})
	}

	return &output, nil
}

func GetOneCertificate(context *base.ApplicationContext, studentPrivateKeyString string, certUUID string) (*data.CertificateOutputFull, error) {
	result, err := educert.QueryCertificateByUUID(context.Contract, certUUID)
	if err != nil {
		return nil, err
	}

	var cert data.BlockchainCertificateOutput
	err = json.Unmarshal([]byte(result), &cert)
	if err != nil {
		return nil, err
	}

	row := context.DB.QueryRow("SELECT student_encrypted_ks_1, student_encrypted_ks_2, ks_encrypted_file FROM certificates WHERE uuid = ?", certUUID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var studentEncryptedKS1, studentEncryptedKS2, ksEncryptedFile string
	err = row.Scan(&studentEncryptedKS1, &studentEncryptedKS2, &ksEncryptedFile)
	if err != nil {
		return nil, err
	}

	ks, err := utils.ElGamalDecryptString(studentPrivateKeyString, studentEncryptedKS1, studentEncryptedKS2)
	if err != nil {
		return nil, err
	}

	base64File, err := utils.VigenereDecryptString(ks, ksEncryptedFile)
	if err != nil {
		return nil, err
	}

	return &data.CertificateOutputFull{
		Base64File: base64File,
		CertificateOutput: data.CertificateOutput{
			CertHash:            cert.CertHash,
			CertUUID:            cert.CertUUID,
			StudentPublicKey:    cert.StudentPublicKey,
			UniversityPublicKey: cert.UniversityPublicKey,
			DateOfIssuing:       cert.DateOfIssuing,
		},
	}, nil
}
