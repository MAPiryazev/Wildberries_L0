package repository

import (
	"database/sql"
	"errors"
)

func (orderRepo *orderRepo) cityIDFromName(tx *sql.Tx, name string) (int, error) {
	var id int
	err := tx.QueryRow(`SELECT id FROM city WHERE name = $1`, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(`INSERT INTO city(name) VALUES($1) RETURNING id`, name).Scan(&id)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (orderRepo *orderRepo) regionIDFromName(tx *sql.Tx, name string) (int, error) {
	var id int
	err := tx.QueryRow(`SELECT id FROM region WHERE name = $1`, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(`INSERT INTO region(name) VALUES($1) RETURNING id`, name).Scan(&id)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (orderRepo *orderRepo) currencyIDFromName(tx *sql.Tx, name string) (int, error) {
	var id int
	err := tx.QueryRow(`SELECT id FROM currency WHERE name = $1`, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(`INSERT INTO currency(name) VALUES($1) RETURNING id`, name).Scan(&id)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (orderRepo *orderRepo) brandIDFromName(tx *sql.Tx, name string) (int, error) {
	var id int
	err := tx.QueryRow(`SELECT id FROM brand WHERE name = $1`, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(`INSERT INTO brand(name) VALUES($1) RETURNING id`, name).Scan(&id)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (orderRepo *orderRepo) bankIDFromName(tx *sql.Tx, name string) (int, error) {
	var id int
	err := tx.QueryRow(`SELECT id FROM bank WHERE name = $1`, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(`INSERT INTO bank(name) VALUES($1) RETURNING id`, name).Scan(&id)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (orderRepo *orderRepo) deliveryServiceIDFromName(tx *sql.Tx, name string) (int, error) {
	var id int
	err := tx.QueryRow(`SELECT id FROM delivery_service WHERE name = $1`, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(`INSERT INTO delivery_service(name) VALUES($1) RETURNING id`, name).Scan(&id)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
