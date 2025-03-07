package universities

import (
	"certiblock/base"
	"certiblock/base/data"
	"certiblock/utils"
	"fmt"
)

func GetById(context *base.ApplicationContext, id int) (*data.BCUniversity, error) {
	row := context.DB.QueryRow("SELECT `id`, `name`, `public_key` FROM `bc_universities` WHERE `id` = ?;", id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	university := data.BCUniversity{}
	err := row.Scan(&university.ID, &university.Name, &university.PublicKey)
	if err != nil {
		return nil, err
	}

	return &university, nil
}

func GetByPrivateKey(context *base.ApplicationContext, privateKey string) (*data.BCUniversity, error) {
	fmt.Printf("Private key = %v\n", privateKey)
	publicKey := utils.PrivateKeyToPublicKey(privateKey)
	fmt.Printf("Public key = %v\n", publicKey)

	row := context.DB.QueryRow("SELECT `id`, `name`, `public_key` FROM `bc_universities` WHERE `public_key` = ?;", publicKey)
	if row.Err() != nil {
		return nil, row.Err()
	}

	university := data.BCUniversity{}
	err := row.Scan(&university.ID, &university.Name, &university.PublicKey)
	if err != nil {
		return nil, err
	}

	return &university, nil
}

func Register(context *base.ApplicationContext, universityInput data.BCUniversityInput) (*data.BCUniversityOutput, error) {
	if universityInput.AdminApprovalKey != "we approved this university to join us" {
		return nil, fmt.Errorf("admin approval key is invalid")
	}
	s, err := GetByPrivateKey(context, universityInput.PrivateKey)
	if s != nil {
		return nil, fmt.Errorf("key already exists, generate another")
	}
	// if err != nil {
	// 	return nil, err
	// }

	publicKey := utils.PrivateKeyToPublicKey(universityInput.PrivateKey)
	result, err := context.DB.Exec("INSERT INTO `bc_universities`(name, public_key) VALUES(?, ?);", universityInput.Name, publicKey)
	if err != nil {
		return nil, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	output := data.BCUniversityOutputResponse(data.BCUniversityOutput{
		BCUniversity: data.BCUniversity{
			BCUniversityCommon: data.BCUniversityCommon{
				Name: universityInput.Name,
			},
			ID:        int(lastInsertedID),
			PublicKey: publicKey,
		},
	})
	return &output, nil
}
