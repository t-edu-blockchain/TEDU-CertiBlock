package countries

import (
	"certiblock/base"
	"certiblock/base/data"
)

func GetById(context *base.ApplicationContext, id int) (*data.Country, error) {
	// query country
	row := context.DB.QueryRow("SELECT `id`, `name` FROM `countries` WHERE `id` = ?;", id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	country := data.Country{}
	err := row.Scan(&country.ID, &country.Name)
	if err != nil {
		return nil, err
	}

	// return country
	return &country, nil
}

func GetAll(context *base.ApplicationContext) ([]data.Country, error) {
	// query countries
	rows, err := context.DB.Query("SELECT `id`, `name` FROM `countries`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over rows
	countries := []data.Country{}
	for rows.Next() {
		country := data.Country{}
		err := rows.Scan(&country.ID, &country.Name)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	// return countries
	return countries, nil
}
