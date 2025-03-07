package students

import (
	"certiblock/base"
	"certiblock/base/data"
	"fmt"
)

func GetByPublicKey(context *base.ApplicationContext, public_key string) (*data.BCStudent, error) {
	// query student
	row := context.DB.QueryRow("SELECT `id`, `public_key` FROM `bc_students` WHERE `public_key` = ?;", public_key)
	if row.Err() != nil {
		return nil, row.Err()
	}

	student := data.BCStudent{}
	err := row.Scan(&student.ID, &student.PublicKey)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func Register(context *base.ApplicationContext, studentInput data.BCStudentInput) (*data.BCStudentOutput, error) {
	// check if student already exists
	s, err := GetByPublicKey(context, studentInput.PublicKey)
	if s != nil {
		return nil, fmt.Errorf("public key already exists, generate another")
	}
	if err != nil {
		return nil, err
	}

	// insert student
	result, err := context.DB.Exec("INSERT INTO `bc_students` (public_key) VALUES (?);", studentInput.PublicKey)
	if err != nil {
		return nil, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	studentOutput := data.BCStudentOutputResponse(data.BCStudent{
		ID:             int(lastInsertedID),
		BCStudentInput: studentInput,
	})

	return &studentOutput, nil
}
