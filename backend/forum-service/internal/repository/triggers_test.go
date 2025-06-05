package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTriggers_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("CREATE OR REPLACE FUNCTION update_comment_count").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TRIGGER IF EXISTS update_comment_count_insert").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TRIGGER IF EXISTS update_comment_count_delete").WillReturnResult(sqlmock.NewResult(0, 0))

	err = CreateTriggers(db)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTriggers_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("CREATE OR REPLACE FUNCTION update_comment_count").WillReturnError(errors.New("fail"))

	err = CreateTriggers(db)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
