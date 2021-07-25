package dao

import (
	"database/sql"
	"testing"

	"github.com/pkg/errors"
)

func TestGetOneRow(t *testing.T) {
	id := "1"
	_, err := GetOneRow(id)
	if err != nil{
		if errors.Cause(err) == sql.ErrNoRows{
			t.Errorf("can't get row from database with id: %s, err: %+v", id, err)
			return
		}else{
			t.Errorf("unknown error: %+v", err)
			return
		}
	}
	t.Log("get one row succeed with id: ", id)
}
