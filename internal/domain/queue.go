package domain

import "database/sql"

type Queue struct {
	ID      int
	UserID  int
	VideoID sql.NullInt16
	Folder  string
	Title   string
	Idx     sql.NullInt16
	Type    string
}
