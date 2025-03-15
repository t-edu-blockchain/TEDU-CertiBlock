package qrs

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/base/data"
	educert "CertiBlock/application/backend/gateway"
	"CertiBlock/application/shared/utils"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// {
// 	"certHash": "999999999",
// 	"universitySignature": "123456789",
// 	"studentSignature": "987654321",
// 	"dateOfIssuing": "2025-03-08",
// 	"certNumber": "",
// 	"certUUID": "111111111",
// 	"universityPK": "222222222",
// 	"studentPK": "333333333",
// 	"dataType": "certificate"
//   }

func CreateQR(context *base.ApplicationContext, qrInput *data.QRInput) (*data.QROutput, error) {
	result, err := educert.QueryCertificateByUUID(context.Contract, qrInput.CertUUID)
	if err != nil {
		return nil, err
	}

	var cert data.BlockchainCertificateOutput
	err = json.Unmarshal([]byte(result), &cert)
	if err != nil {
		return nil, err
	}

	qrPrivateKeyString, err := utils.GeneratePrivateKeyString()
	if err != nil {
		return nil, err
	}

	qrPublicKeyString, err := utils.ComputePublicKeyString(qrPrivateKeyString)
	if err != nil {
		return nil, err
	}

	rkString, err := utils.PREComputeReEncryptionKeyString(qrInput.StudentPrivateKey, qrPublicKeyString)
	if err != nil {
		return nil, err
	}

	row := context.DB.QueryRow("SELECT student_encrypted_ks_1, student_encrypted_ks_2, ks_encrypted_file FROM students WHERE uuid = ?", qrInput.CertUUID)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var studentEncryptedKS1, studentEncryptedKS2, ksEncryptedFile string
	err = row.Scan(&studentEncryptedKS1, &studentEncryptedKS2, &ksEncryptedFile)
	if err != nil {
		return nil, err
	}

	qrEncryptedKS1, qrEncryptedKS2, err := utils.PREEncryptString(rkString, studentEncryptedKS1, studentEncryptedKS2)
	if err != nil {
		return nil, err
	}

	qrInDb, err := context.DB.Exec("INSERT INTO qrs (cert_uuid, qr_encrypted_ks_1, qr_encrypted_ks_2) VALUES (?, ?, ?)", qrInput.CertUUID, qrEncryptedKS1, qrEncryptedKS2)
	if err != nil {
		return nil, err
	}

	qrID, err := qrInDb.LastInsertId()
	if err != nil {
		return nil, err
	}

	qrOutput := data.QROutput{
		Content: strconv.FormatInt(qrID, 10) + "0a0" + qrPrivateKeyString,
	}

	return &qrOutput, nil
}

func VerifyAndDecodeQR(context *base.ApplicationContext, token string) (*string, error) {
	// parse token
	tokenParts := strings.Split(token, "0a0")
	if len(tokenParts) != 2 {
		return nil, errors.New("invalid token format")
	}

	qrID, err := strconv.Atoi(tokenParts[0])
	if err != nil {
		return nil, err
	}

	qrPrivateKeyString := tokenParts[1]

	// query qr
	row := context.DB.QueryRow("SELECT cert_uuid, qr_encrypted_ks_1, qr_encrypted_ks_2 FROM qrs WHERE id = ?", qrID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var certUUID, qrEncryptedKS1, qrEncryptedKS2 string
	err = row.Scan(&certUUID, &qrEncryptedKS1, &qrEncryptedKS2)
	if err != nil {
		return nil, err
	}

	// get the cert from the blockchain
	rawCertInBlockchain, err := educert.QueryCertificateByUUID(context.Contract, certUUID)
	if err != nil {
		return nil, err
	}
	var certInBlockchain data.BlockchainCertificateOutput
	err = json.Unmarshal([]byte(rawCertInBlockchain), &certInBlockchain)
	if err != nil {
		return nil, err
	}

	// decrypt ks
	ks, err := utils.PREDecryptString(certInBlockchain.StudentPublicKey, qrPrivateKeyString, qrEncryptedKS1, qrEncryptedKS2)
	if err != nil {
		return nil, err
	}

	// get the file
	row = context.DB.QueryRow("SELECT ks_encrypted_file FROM certificates WHERE uuid = ?", certUUID)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var ksEncryptedFile string
	err = row.Scan(&ksEncryptedFile)
	if err != nil {
		return nil, err
	}

	file, err := utils.VigenereDecryptString(ks, ksEncryptedFile)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
