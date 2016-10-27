package pgsql

import (
	"database/sql"
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
)

type ProfileService struct{}

type profileModel struct {
	id        nullString
	userID    nullString
	firstName nullString
	lastName  nullString
	country   nullString
	city      nullString
	state     nullString
	zipCode   nullString
	address   nullString
	phone     nullString
	company   nullString
	created   nullTime
	updated   nullTime
}

func (pm *profileModel) convert() *model.UserProfile {

	var p = model.UserProfile{}
	p.UUID = pm.id.String
	p.UserID = pm.userID.String
	p.FirstName = pm.firstName.String
	p.LastName = pm.lastName.String
	p.Country = pm.country.String
	p.City = pm.city.String
	p.State = pm.state.String
	p.ZipCode = pm.zipCode.String
	p.Address = pm.address.String
	p.Phone = pm.phone.String
	p.Company = pm.company.String
	p.Created = pm.created.Time
	p.Updated = pm.updated.Time

	return &p
}

func (ProfileService) Insert(db adapter.IDatabase, userID string) (*string, *e.Err) {

	var err error
	var id nullString

	const query = `
		INSERT INTO profiles (user_id)
		VALUES ($1)
		RETURNING id`

	err = db.QueryRow(query, userID).Scan(&id)
	if err != nil {
		return nil, e.Profile.Unknown(err)
	}

	return &id.String, nil
}

func (ProfileService) Update(db adapter.IDatabase, userID, firstName, lastName, company, city string) (bool, *e.Err) {

	var err error

	const sqlquery = `
		UPDATE profiles
		SET
			first_name = $2,
			last_name = $3,
			city = $4,
			company = $5,
			updated = now()
		WHERE user_id = $1`

	res, err := db.Exec(sqlquery, userID, firstName, lastName, city, company)
	if err != nil {
		return false, e.Profile.Unknown(err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return false, e.Profile.Unknown(err)
	}

	return rowCount != 0, nil
}

func (ProfileService) GetByUserID(db adapter.IDatabase, userID string) (*model.UserProfile, *e.Err) {

	var err error
	var pm = new(profileModel)

	const sqlquery = `
		SELECT
			id,
			first_name,
			last_name,
			country,
			city,
			state,
			zip_code,
			address,
			phone,
			company
		FROM profiles
		WHERE user_id = $1`

	err = db.QueryRow(sqlquery, userID).Scan(&pm.id, &pm.firstName, &pm.lastName, &pm.country, &pm.city, &pm.state,
		&pm.zipCode, &pm.address, &pm.phone, &pm.company)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, e.Profile.NotFound(err)
	default:
		return nil, e.Profile.Unknown(err)
	}

	var profile = pm.convert()

	return profile, nil
}
