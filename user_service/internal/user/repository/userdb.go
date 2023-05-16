package repository

import "database/sql"

type user struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Login       string         `db:"login"`
	Password    string         `db:"password"`
	Description sql.NullString `db:"description"`
}
