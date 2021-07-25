package dao

import (
	"database/sql"
	"github.com/pkg/errors"
)

func GetOneRow(id string)(*sql.Row, error){

	// do some query stuff with id

	mockRow := &sql.Row{}
	mockErr := sql.ErrNoRows

	if mockErr != nil{
		return nil, errors.Wrap(mockErr, "pkg/dao/sql.go: line 12: get one row error")
	}
	return mockRow, nil
}