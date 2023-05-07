package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"testing"
	"transcribify/internal/models"
	"transcribify/pkg/dbclient"
)

func TestUserRepository_PutUser(t *testing.T) {
	type Case struct {
		Name        string
		ErrExpected error

		User         models.User
		ExpectedUser models.User
	}

	tc := []Case{
		{
			Name:         "Put user in db",
			ErrExpected:  nil,
			User:         models.User{Email: "test1@gmail.com", Password: "1234567890"},
			ExpectedUser: models.User{ID: 1, Email: "test1@gmail.com", Password: "1234567890"},
		},
		{
			Name:         "Put user with empty password",
			ErrExpected:  &pgconn.PgError{Severity: "ERROR", Code: "22000", Message: "empty email or password", Detail: "", Hint: "enter user email and password", Position: 0, InternalPosition: 0, InternalQuery: "", Where: "PL/pgSQL function put_user(text,text) line 5 at RAISE", SchemaName: "", TableName: "", ColumnName: "", DataTypeName: "", ConstraintName: "", File: "pl_exec.c", Line: 3893, Routine: "exec_stmt_raise"},
			User:         models.User{Email: "test2@gmail.com", Password: ""},
			ExpectedUser: models.User{Email: "test2@gmail.com", Password: ""},
		},
		{
			Name:         "Put user with empty login",
			ErrExpected:  &pgconn.PgError{Severity: "ERROR", Code: "22000", Message: "empty email or password", Detail: "", Hint: "enter user email and password", Position: 0, InternalPosition: 0, InternalQuery: "", Where: "PL/pgSQL function put_user(text,text) line 5 at RAISE", SchemaName: "", TableName: "", ColumnName: "", DataTypeName: "", ConstraintName: "", File: "pl_exec.c", Line: 3893, Routine: "exec_stmt_raise"},
			User:         models.User{Email: "", Password: "1234567890"},
			ExpectedUser: models.User{Email: "", Password: "1234567890"},
		},
		{
			Name:         "Put user with empty login and empty password",
			ErrExpected:  &pgconn.PgError{Severity: "ERROR", Code: "22000", Message: "empty email or password", Detail: "", Hint: "enter user email and password", Position: 0, InternalPosition: 0, InternalQuery: "", Where: "PL/pgSQL function put_user(text,text) line 5 at RAISE", SchemaName: "", TableName: "", ColumnName: "", DataTypeName: "", ConstraintName: "", File: "pl_exec.c", Line: 3893, Routine: "exec_stmt_raise"},
			User:         models.User{Email: "", Password: ""},
			ExpectedUser: models.User{Email: "", Password: ""},
		},
	}

	ctx := context.Background()
	db, err := dbclient.NewClient(ctx)
	if err != nil {
		t.Error(err)
	}
	defer db.Close(ctx)
	repo := NewUserRepository(db)
	for _, c := range tc {
		t.Run(c.Name, func(t *testing.T) {

			user := &c.User
			err = repo.PutUser(ctx, user)
			assert.Equal(t, c.ErrExpected, err)
			assert.Equal(t, c.ExpectedUser, *user)
		})
	}
}

func TestUserRepository_GetUserByLoginPassword(t *testing.T) {
	type Case struct {
		Name        string
		ErrExpected error

		User         models.User
		ExpectedUser models.User
	}

	tc := []Case{
		{
			Name:         "Get user by login and password",
			ErrExpected:  nil,
			User:         models.User{Email: "test1@gmail.com", Password: "1234567890"},
			ExpectedUser: models.User{ID: 1, Email: "test1@gmail.com", Password: "1234567890"},
		},
		{
			Name:         "Get user by login",
			ErrExpected:  nil,
			User:         models.User{Email: "test1@gmail.com", Password: ""},
			ExpectedUser: models.User{ID: 1, Email: "test1@gmail.com", Password: "1234567890"},
		},
		{
			Name:         "Get user by password",
			ErrExpected:  &pgconn.PgError{Severity: "ERROR", Code: "22000", Message: "enter login", Detail: "", Hint: "cannot fund user with empty login", Position: 0, InternalPosition: 0, InternalQuery: "", Where: "PL/pgSQL function get_user(text,text) line 5 at RAISE", SchemaName: "", TableName: "", ColumnName: "", DataTypeName: "", ConstraintName: "", File: "pl_exec.c", Line: 3893, Routine: "exec_stmt_raise"},
			User:         models.User{Email: "", Password: "1234567890"},
			ExpectedUser: models.User{Email: "", Password: "1234567890"},
		},
		{
			Name:         "Get user with empty login and password",
			ErrExpected:  &pgconn.PgError{Severity: "ERROR", Code: "22000", Message: "enter login", Detail: "", Hint: "cannot fund user with empty login", Position: 0, InternalPosition: 0, InternalQuery: "", Where: "PL/pgSQL function get_user(text,text) line 5 at RAISE", SchemaName: "", TableName: "", ColumnName: "", DataTypeName: "", ConstraintName: "", File: "pl_exec.c", Line: 3893, Routine: "exec_stmt_raise"},
			User:         models.User{Email: "", Password: ""},
			ExpectedUser: models.User{Email: "", Password: ""},
		},
	}
	ctx := context.Background()
	// connect with db
	db, err := dbclient.NewClient(ctx)
	if err != nil {
		t.Error(err)
	}
	defer db.Close(ctx)
	repo := NewUserRepository(db)

	// put example to get
	err = repo.PutUser(ctx, &tc[0].User)
	if err != nil {
		t.Error(err)
	}
	for _, c := range tc {
		t.Run(c.Name, func(t *testing.T) {

			user := &c.User
			err = repo.GetUserByLoginPassword(ctx, user)

			assert.Equal(t, c.ErrExpected, err)
			assert.Equal(t, c.ExpectedUser, *user)
		})
	}
}

func TestUserRepository_PutUserVideo(t *testing.T) {
	type UserVideo struct {
		UID int
		VID string
	}
	type Case struct {
		UV          UserVideo
		Name        string
		ErrExpected error
	}
	tc := []Case{
		{
			Name:        "Adding user video",
			UV:          UserVideo{UID: 1, VID: "00000000000"},
			ErrExpected: nil,
		},
		{
			Name:        "Adding user video without UID",
			UV:          UserVideo{VID: "00000000000"},
			ErrExpected: &pgconn.PgError{Severity: "ERROR", Code: "22000", Message: "empty user-id or video-id", Detail: "", Hint: "enter user-id or video-id", Position: 0, InternalPosition: 0, InternalQuery: "", Where: "PL/pgSQL function put_user_video(integer,character) line 5 at RAISE", SchemaName: "", TableName: "", ColumnName: "", DataTypeName: "", ConstraintName: "", File: "pl_exec.c", Line: 3893, Routine: "exec_stmt_raise"},
		},
		{
			Name:        "Adding user video without VID",
			UV:          UserVideo{UID: 1},
			ErrExpected: &pgconn.PgError{Severity: "ERROR", Code: "22000", Message: "empty user-id or video-id", Detail: "", Hint: "enter user-id or video-id", Position: 0, InternalPosition: 0, InternalQuery: "", Where: "PL/pgSQL function put_user_video(integer,character) line 5 at RAISE", SchemaName: "", TableName: "", ColumnName: "", DataTypeName: "", ConstraintName: "", File: "pl_exec.c", Line: 3893, Routine: "exec_stmt_raise"},
		},
		{
			Name:        "Adding user video without UID and VID",
			UV:          UserVideo{},
			ErrExpected: &pgconn.PgError{Severity: "ERROR", Code: "22000", Message: "empty user-id or video-id", Detail: "", Hint: "enter user-id or video-id", Position: 0, InternalPosition: 0, InternalQuery: "", Where: "PL/pgSQL function put_user_video(integer,character) line 5 at RAISE", SchemaName: "", TableName: "", ColumnName: "", DataTypeName: "", ConstraintName: "", File: "pl_exec.c", Line: 3893, Routine: "exec_stmt_raise"},
		},
	}

	ctx := context.Background()
	client, err := dbclient.NewClient(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	repo := NewUserRepository(client)
	//set-up user and video
	v := NewYTVideoRepository(client)
	video := new(models.YTVideo)
	video.Title = "Title"
	_, err = v.CreateVideo(ctx, models.VideoRequest{VideoID: "00000000000", Language: "ua"}, video)
	if err != nil {
		t.Error(err)
	}

	err = repo.PutUser(ctx, &models.User{Email: "test@gmail.com", Password: "123456789"})
	if err != nil {
		t.Error(err)
	}

	for _, c := range tc {
		t.Run(c.Name, func(t *testing.T) {

			err = repo.PutUserVideo(ctx, c.UV.UID, c.UV.VID)
			assert.Equal(t, c.ErrExpected, err)
		})
	}
}

//func TestUserRepository_GetUserVideos(t *testing.T) {
//
//}
