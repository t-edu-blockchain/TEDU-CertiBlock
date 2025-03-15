package students

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/base/data"
	"CertiBlock/application/backend/certiblock/services/countries"
	"CertiBlock/application/shared/utils"

	"errors"
)

func RegisterStudent(context *base.ApplicationContext, studentInput *data.StudentInput) (*data.StudentOutput, error) {
	country, err := countries.GetById(context, studentInput.CountryID)
	if err != nil || country == nil {
		return nil, errors.New("country not found")
	}

	privateKeyString := utils.HashSHA512(
		country.Name + studentInput.NIN + studentInput.FullName + studentInput.DateOfBirth + studentInput.Password,
	)
	publicKeyString, err := utils.ComputePublicKeyString(privateKeyString)
	if err != nil {
		return nil, errors.New("error computing public key")
	}

	row := context.DB.QueryRow("SELECT public_key FROM students WHERE public_key = ?", publicKeyString)
	var existingPublicKey string
	err = row.Scan(&existingPublicKey)
	if err == nil {
		return nil, errors.New("student already registered")
	}

	encryptedPartialPersonalInfo1, encryptedPartialPersonalInfo2, err := utils.ElGamalEncryptString(
		publicKeyString,
		studentInput.FullName,
	)
	if err != nil {
		return nil, errors.New("error encrypting partial personal info")
	}

	_, err = context.DB.Exec(
		"INSERT INTO students (public_key, encrypted_partial_personal_info1, encrypted_partial_personal_info2) VALUES (?, ?, ?)",
		publicKeyString,
		encryptedPartialPersonalInfo1,
		encryptedPartialPersonalInfo2,
	)
	if err != nil {
		return nil, errors.New("error inserting student")
	}

	return &data.StudentOutput{
		PrivateKey: privateKeyString,
		PublicKey:  publicKeyString,
	}, nil
}

func LoginStudent(context *base.ApplicationContext, studentInput *data.StudentInput) (*data.StudentOutput, error) {
	country, err := countries.GetById(context, studentInput.CountryID)
	if err != nil || country == nil {
		return nil, errors.New("country not found")
	}

	privateKeyString := utils.HashSHA512(
		country.Name + studentInput.NIN + studentInput.FullName + studentInput.DateOfBirth,
	)
	publicKeyString, err := utils.ComputePublicKeyString(privateKeyString)
	if err != nil {
		return nil, errors.New("error computing public key")
	}

	row := context.DB.QueryRow("SELECT public_key FROM students WHERE public_key = ?", publicKeyString)
	var existingPublicKey string
	err = row.Scan(&existingPublicKey)
	if err != nil {
		return nil, errors.New("student not registered yet or the password is incorrect")
	}

	return &data.StudentOutput{
		PrivateKey: privateKeyString,
		PublicKey:  publicKeyString,
	}, nil
}
