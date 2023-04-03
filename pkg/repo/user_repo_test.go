package repo

import (
	"context"
	"errors"
	"reddit_backend/pkg/logging"
	"reddit_backend/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		testUser    *User
		prepMockSql func(mock sqlmock.Sqlmock, user *User)
		wantErr     string
	}{
		{ //ok query and exec
			testUser: &User{Id: "123", Username: "test1", password: "12345678"},
			prepMockSql: func(mock sqlmock.Sqlmock, user *User) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"})
				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE`).WithArgs(user.Username).WillReturnRows(rows)
				mock.ExpectExec(`INSERT INTO users`).WithArgs(user.Id, user.Username, user.password).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{ //bad query
			testUser: &User{Id: "123", Username: "test1", password: "12345678"},
			prepMockSql: func(mock sqlmock.Sqlmock, user *User) {
				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE`).WithArgs(user.Username).WillReturnError(errors.New("bad query"))
			},
			wantErr: "bad query",
		},
		{ //user exist
			testUser: &User{Id: "123", Username: "test1", password: "12345678"},
			prepMockSql: func(mock sqlmock.Sqlmock, user *User) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(user.Id, user.Username, user.password)
				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE`).WithArgs(user.Username).WillReturnRows(rows)
			},
			wantErr: "username already exists",
		},
		{ //bad exec
			testUser: &User{Id: "123", Username: "test1", password: "12345678"},
			prepMockSql: func(mock sqlmock.Sqlmock, user *User) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"})
				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE`).WithArgs(user.Username).WillReturnRows(rows)
				mock.ExpectExec(`INSERT INTO users`).WithArgs(user.Id, user.Username, user.password).WillReturnError(errors.New("bad exec"))
			},
			wantErr: "bad exec",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	zapLogger := zap.NewNop()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	loggerCtx := &logging.Logger{Zap: logger, RequestIDKey: "requestID", LoggerKey: "logger"}
	repoUser := &UserMySql{Logger: loggerCtx, MySql: db}

	for _, tc := range testCases {
		tc.prepMockSql(mock, tc.testUser)
		user, err := repoUser.CreateUser(context.TODO(), tc.testUser.Id, tc.testUser.Username, tc.testUser.password)
		if err != nil {
			if err.Error() != tc.wantErr {
				t.Errorf("expected err: %s, got: %s", tc.wantErr, err)
				continue
			}
			continue
		}
		require.EqualValues(t, tc.testUser.Id, user.Id)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestCheckPasswordUser(t *testing.T) {
	testCases := []struct {
		testUser    *User
		prepMockSql func(mock sqlmock.Sqlmock, user *User)
		wantErr     string
	}{
		{ //ok query and exec
			testUser: &User{Id: "123", Username: "test1", password: "12345678"},
			prepMockSql: func(mock sqlmock.Sqlmock, user *User) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(user.Id, user.Username, user.password)
				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE`).WithArgs(user.Username).WillReturnRows(rows)
			},
		},
		{ //not ok password
			testUser: &User{Id: "123", Username: "test1", password: "12345678"},
			prepMockSql: func(mock sqlmock.Sqlmock, user *User) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(user.Id, user.Username, "1111111111")
				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE`).WithArgs(user.Username).WillReturnRows(rows)
			},
			wantErr: utils.ErrBadPass.Error(),
		},
		{ //user not found
			testUser: &User{Id: "123", Username: "test1", password: "12345678"},
			prepMockSql: func(mock sqlmock.Sqlmock, user *User) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"})
				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE`).WithArgs(user.Username).WillReturnRows(rows)
			},
			wantErr: utils.ErrNoUser.Error(),
		},
		{ //bad query
			testUser: &User{Id: "123", Username: "test1", password: "12345678"},
			prepMockSql: func(mock sqlmock.Sqlmock, user *User) {
				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE`).WithArgs(user.Username).WillReturnError(errors.New("bad query"))
			},
			wantErr: "bad query",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	zapLogger := zap.NewNop()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	loggerCtx := &logging.Logger{Zap: logger, RequestIDKey: "requestID", LoggerKey: "logger"}
	repoUser := &UserMySql{Logger: loggerCtx, MySql: db}

	for _, tc := range testCases {
		tc.prepMockSql(mock, tc.testUser)
		user, err := repoUser.CheckPasswordUser(context.TODO(), tc.testUser.Username, tc.testUser.password)
		if err != nil {
			if err.Error() != tc.wantErr {
				t.Errorf("expected err: %s, got: %s", tc.wantErr, err)
				continue
			}
			continue
		}
		require.EqualValues(t, tc.testUser.Id, user.Id)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
