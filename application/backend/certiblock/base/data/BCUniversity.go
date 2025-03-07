package data

type BCUniversityLogin struct {
	PrivateKey string `json:"private_key"`
}

type BCUniversityCommon struct {
	Name string `json:"name"`
}

type BCUniversityInput struct {
	BCUniversityCommon
	AdminApprovalKey string `json:"admin_approval_key"`
	PrivateKey       string `json:"private_key"`
}

type BCUniversity struct {
	BCUniversityCommon
	ID        int    `json:"id"`
	PublicKey string `json:"public_key"`
}

type BCUniversityOutput struct {
	BCUniversity
}

func BCUniversityOutputResponse(bcUniversity any) BCUniversityOutput {
	switch c := bcUniversity.(type) {
	case BCUniversity:
		return BCUniversityOutput{
			BCUniversity: c,
		}
	case *BCUniversity:
		return BCUniversityOutput{
			BCUniversity: *c,
		}
	case BCUniversityOutput:
		return c
	case *BCUniversityOutput:
		return *c
	default:
		// Return an empty BCUniversityOutput if conversion is not possible
		return BCUniversityOutput{}
	}
}
