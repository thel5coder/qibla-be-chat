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
	FindByID(id int) (UserEntity, error)
}

// UserEntity ....
type UserEntity struct {
	ID           int            `db:"id"`
	CompanyID    sql.NullInt64  `db:"companyId"`
	RoleID       sql.NullInt64  `db:"roleId"`
	Name         sql.NullString `db:"name"`
	Email        sql.NullString `db:"email"`
	EmailValidAt sql.NullString `db:"emailValidAt"`
	Phone        sql.NullString `db:"phone"`
	PhoneValidAt sql.NullString `db:"phoneValidAt"`
	Password     sql.NullString `db:"password"`
	Photo        sql.NullString `db:"photo"`
	Status       sql.NullBool   `db:"status"`
	CreatedAt    sql.NullString `db:"createdAt"`
	UpdatedAt    sql.NullString `db:"updatedAt"`
	DeletedAt    sql.NullString `db:"deletedAt"`
}

// NewUserModel ...
func NewUserModel(db *sql.DB) IUser {
	return &userModel{DB: db}
}

// FindByID ...
func (model userModel) FindByID(id int) (data UserEntity, err error) {
	query :=
		`SELECT u."id", u."companyId", u."roleId", u."name", u."email", u."emailValidAt", u."phone",
		u."phoneValidAt", u."password", u."status", u."photo", u."createdAt", u."updatedAt", u."deletedAt"
		FROM "Users" u
		WHERE u."deletedAt" IS NULL AND u."id" = $1
		ORDER BY u."createdAt" DESC LIMIT 1`
	err = model.DB.QueryRow(query, id).Scan(
		&data.ID, &data.CompanyID, &data.RoleID, &data.Name, &data.Email, &data.EmailValidAt, &data.Phone,
		&data.PhoneValidAt, &data.Password, &data.Status, &data.Photo, &data.CreatedAt, &data.UpdatedAt,
		&data.DeletedAt,
	)

	return data, err
}
