package model

import (
	"database/sql"
)

// userModel ...
type userModel struct {
	DB *sql.DB
}

// IUser ...
type IUser interface {
	FindByID(id string) (UserEntity, error)
	FindByOdooUserID(odooUserID int64) (UserEntity, error)
}

// UserEntity ....
type UserEntity struct {
	ID        string         `db:"id"`
	Username  sql.NullString `db:"username"`
	Email     sql.NullString `db:"email"`
	Name      sql.NullString `db:"name"`
	Password  sql.NullString `db:"password"`
	RoleID    sql.NullString `db:"role_id"`
	RoleName  sql.NullString `db:"role_name"`
	OdoUserID sql.NullInt64  `db:"odo_user_id"`
	IsActive  sql.NullBool   `db:"is_active"`
	CreatedAt sql.NullString `db:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// NewUserModel ...
func NewUserModel(db *sql.DB) IUser {
	return &userModel{DB: db}
}

// FindByID ...
func (model userModel) FindByID(id string) (data UserEntity, err error) {
	query :=
		`SELECT u."id", u."username", u."email", u."name", u."password", u."role_id", r."name", u."odo_user_id",
		u."is_active", u."created_at", u."updated_at", u."deleted_at"
		FROM "users" u
		LEFT JOIN "roles" r ON r."id" = u."role_id"
		WHERE u."deleted_at" IS NULL AND u."id" = $1
		ORDER BY u."created_at" DESC LIMIT 1`
	err = model.DB.QueryRow(query, id).Scan(
		&data.ID, &data.Username, &data.Email, &data.Name, &data.Password, &data.RoleID, &data.RoleName,
		&data.OdoUserID, &data.IsActive, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
	)

	return data, err
}

// FindByOdooUserID ...
func (model userModel) FindByOdooUserID(odooUserID int64) (data UserEntity, err error) {
	query :=
		`SELECT u."id", u."username", u."email", u."name", u."password", u."role_id", r."name", u."odo_user_id",
		u."is_active", u."created_at", u."updated_at", u."deleted_at"
		FROM "users" u
		LEFT JOIN "roles" r ON r."id" = u."role_id"
		WHERE u."deleted_at" IS NULL AND u."odo_user_id" = $1
		ORDER BY u."created_at" DESC LIMIT 1`
	err = model.DB.QueryRow(query, odooUserID).Scan(
		&data.ID, &data.Username, &data.Email, &data.Name, &data.Password, &data.RoleID, &data.RoleName,
		&data.OdoUserID, &data.IsActive, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
	)

	return data, err
}
