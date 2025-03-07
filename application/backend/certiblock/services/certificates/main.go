package certificates

import (
	"certiblock/base"
	"certiblock/base/data"
	"certiblock/services/universities"
	"certiblock/utils"
)

func private_BCCertificateInputRequest(bcCertificateInput data.BCCertificateInput) (*data.BCCertificateInputValidated, error) {
	h := utils.HashSHA512(bcCertificateInput.File)
	K_s := utils.GenerateSecureRandomString(128)

	universityPublicKey := utils.PrivateKeyToPublicKey(bcCertificateInput.UniversityPrivateKey)

	C_X := utils.Encrypt(K_s, bcCertificateInput.File)

	C_U := utils.Encrypt(universityPublicKey, K_s)
	C_S := utils.Encrypt(bcCertificateInput.StudentPublicKey, K_s)

	S_0 := utils.Sign(bcCertificateInput.UniversityPrivateKey, h+C_X)
	S := utils.Sign(bcCertificateInput.UniversityPrivateKey, h+C_U+C_S+bcCertificateInput.StudentPublicKey+S_0)

	return &data.BCCertificateInputValidated{
		V_h:                  h,
		V_C_X:                C_X,
		V_C_U:                C_U,
		V_C_S:                C_S,
		StudentPublicKey:     bcCertificateInput.StudentPublicKey,
		UniversityPrivateKey: bcCertificateInput.UniversityPrivateKey,
		V_S0:                 S_0,
		V_S:                  S,
	}, nil
}

func Issue(context *base.ApplicationContext, certificateInput data.BCCertificateInput) (*data.BCCertificateOutput, error) {
	university, err := universities.GetByPrivateKey(context, certificateInput.UniversityPrivateKey)
	if err != nil {
		return nil, err
	}

	validatedInput, err := private_BCCertificateInputRequest(certificateInput)
	if err != nil {
		return nil, err
	}

	result, err := context.DB.Exec("INSERT INTO bc_certificates (h, C_X, C_U, C_S, student_public_key, university_public_key, S0, S) VALUES (?, ?, ?, ?, ?, ?, ?)",
		validatedInput.V_h,
		validatedInput.V_C_X,
		validatedInput.V_C_U,
		validatedInput.V_C_S,
		validatedInput.StudentPublicKey,
		university.PublicKey,
		validatedInput.V_S0,
		validatedInput.V_S,
	)

	if err != nil {
		return nil, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	output := data.BCCertificateOutputResponse(data.BCCertificate{
		ID:                          int(lastInsertedID),
		BCCertificateInputValidated: *validatedInput,
	})
	return &output, nil
}

func GetCXByHash(context *base.ApplicationContext, hash string) (*string, error) {
	row := context.DB.QueryRow("SELECT C_X FROM bc_certificates WHERE h = ?;", hash)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var C_X string
	err := row.Scan(&C_X)
	if err != nil {
		return nil, err
	}

	return &C_X, nil
}
