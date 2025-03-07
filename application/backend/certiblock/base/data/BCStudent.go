package data

type BCStudentInput struct {
	PublicKey string `json:"public_key"`
}

type BCStudent struct {
	ID int `json:"id"`
	BCStudentInput
}

type BCStudentOutput struct {
	BCStudent
}

func BCStudentOutputResponse(bcStudent any) BCStudentOutput {
	switch c := bcStudent.(type) {
	case BCStudent:
		return BCStudentOutput{
			BCStudent: c,
		}
	case *BCStudent:
		return BCStudentOutput{
			BCStudent: *c,
		}
	case BCStudentOutput:
		return c
	case *BCStudentOutput:
		return *c
	default:
		// Return an empty BCStudentOutput if conversion is not possible
		return BCStudentOutput{}
	}
}
