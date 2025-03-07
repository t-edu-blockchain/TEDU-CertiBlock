package qrs

import (
	"certiblock/base"
	"certiblock/base/data"
	"certiblock/services/certificates"
	"certiblock/utils"
	"fmt"
	"strconv"
)

func CreateQR(context *base.ApplicationContext, qrInput *data.QRInput) (*data.QROutput, error) {
	// For now V_rk is temporarily the student private key itself.

	result, err := context.DB.Exec("INSERT INTO `qrs` (rk, h) VALUES (?, ?);",
		qrInput.V_rk,
		qrInput.V_h,
	)

	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	str := fmt.Sprintf("%v", lastId)
	return &data.QROutput{
		Token: str,
	}, nil
}

func ValidateQR(context *base.ApplicationContext, token string) (*string, error) {
	id, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return nil, err
	}

	row := context.DB.QueryRow("SELECT (rk, h, student_public_key) FROM qrs WHERE id = ?;", id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var rk string
	var h string
	var StudentPublicKey string
	err = row.Scan(&rk, &h, &StudentPublicKey)
	if err != nil {
		return nil, err
	}

	C_X, err := certificates.GetCXByHash(context, h)
	if err != nil {
		return nil, err
	}
	X := utils.Decrypt(rk, *C_X)
	return &X, nil
}
