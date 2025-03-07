package data

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CountryOutput struct {
	Country
}

func CountryOutputResponse(country any) CountryOutput {
	switch c := country.(type) {
	case Country:
		return CountryOutput{
			Country: c,
		}
	case *Country:
		return CountryOutput{
			Country: *c,
		}
	case CountryOutput:
		return c
	case *CountryOutput:
		return *c
	default:
		// Return an empty CountryOutput if conversion is not possible
		return CountryOutput{}
	}
}
