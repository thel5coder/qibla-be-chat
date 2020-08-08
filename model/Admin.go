package model

import (
	"qibla-backend-chat/usecase/viewmodel"
	"database/sql"
	"strings"
	"time"
)

// adminModel ...
type adminModel struct {
	DB *sql.DB
}

// IAdmin ...
type IAdmin interface {
	FindAll(offset, limit int) ([]AdminEntity, int, error)
	FindByID(id string) (AdminEntity, error)
	FindByCode(code string) (AdminEntity, error)
	FindByEmail(email string) (AdminEntity, error)
	Store(body viewmodel.AdminVM, changedAt time.Time) (string, error)
	Update(id string, body viewmodel.AdminVM, changedAt time.Time) (string, error)
	Destroy(id string, changedAt time.Time) (string, error)
}

// AdminEntity ....
type AdminEntity struct {
	ID        string         `db:"id"`
	Code      sql.NullString `db:"code"`
	Name      sql.NullString `db:"name"`
	Email     sql.NullString `db:"email"`
	Password  sql.NullString `db:"password"`
	RoleID    sql.NullString `db:"role_id"`
	Role      RoleEntity     `db:"role"`
	Status    sql.NullBool   `db:"status"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// NewAdminModel ...
func NewAdminModel(db *sql.DB) IAdmin {
	return &adminModel{DB: db}
}

// FindAll ...
func (model adminModel) FindAll(offset, limit int) (data []AdminEntity, count int, err error) {
	query := `SELECT a."id", a."code", a."name", a."email", a."password", a."role_id", a."status",
	a."created_at", a."updated_at", a."deleted_at", r."code", r."name"
	FROM "admins" a
	LEFT JOIN "roles" r ON r."id" = a."role_id"
	WHERE a."deleted_at" IS NULL ORDER BY a."created_at" DESC OFFSET $1 LIMIT $2`
	rows, err := model.DB.Query(query, offset, limit)
	if err != nil {
		return data, count, err
	}

	defer rows.Close()
	for rows.Next() {
		d := AdminEntity{}
		err = rows.Scan(
			&d.ID, &d.Code, &d.Name, &d.Email, &d.Password, &d.RoleID, &d.Status, &d.CreatedAt,
			&d.UpdatedAt, &d.DeletedAt, &d.Role.Code, &d.Role.Name,
		)
		if err != nil {
			return data, count, err
		}
		data = append(data, d)
	}

	err = rows.Err()
	if err != nil {
		return data, count, err
	}

	query = `SELECT COUNT("id") FROM "admins" WHERE "deleted_at" IS NULL`
	err = model.DB.QueryRow(query).Scan(&count)

	return data, count, err
}

// FindByID ...
func (model adminModel) FindByID(id string) (data AdminEntity, err error) {
	query :=
		`SELECT a."id", a."code", a."name", a."email", a."password", a."role_id", a."status",
		a."created_at", a."updated_at", a."deleted_at", r."code", r."name"
		FROM "admins" a
		LEFT JOIN "roles" r ON r."id" = a."role_id"
		WHERE a."deleted_at" IS NULL AND a."id" = $1
		ORDER BY a."created_at" DESC LIMIT 1`
	err = model.DB.QueryRow(query, id).Scan(
		&data.ID, &data.Code, &data.Name, &data.Email, &data.Password, &data.RoleID, &data.Status,
		&data.CreatedAt, &data.UpdatedAt, &data.DeletedAt, &data.Role.Code, &data.Role.Name,
	)

	return data, err
}

// FindByCode ...
func (model adminModel) FindByCode(code string) (data AdminEntity, err error) {
	query :=
		`SELECT a."id", a."code", a."name", a."email", a."password", a."role_id", a."status",
		a."created_at", a."updated_at", a."deleted_at", r."code", r."name"
		FROM "admins" a
		LEFT JOIN "roles" r ON r."id" = a."role_id"
		WHERE a."deleted_at" IS NULL AND a."code" = $1
		ORDER BY a."created_at" DESC LIMIT 1`
	err = model.DB.QueryRow(query, code).Scan(
		&data.ID, &data.Code, &data.Name, &data.Email, &data.Password, &data.RoleID, &data.Status,
		&data.CreatedAt, &data.UpdatedAt, &data.DeletedAt, &data.Role.Code, &data.Role.Name,
	)

	return data, err
}

// FindByEmail ...
func (model adminModel) FindByEmail(email string) (data AdminEntity, err error) {
	query :=
		`SELECT a."id", a."code", a."name", a."email", a."password", a."role_id", a."status",
		a."created_at", a."updated_at", a."deleted_at", r."code", r."name"
		FROM "admins" a
		LEFT JOIN "roles" r ON r."id" = a."role_id"
		WHERE a."deleted_at" IS NULL AND LOWER(a."email") = $1
		ORDER BY a."created_at" DESC LIMIT 1`
	err = model.DB.QueryRow(query, strings.ToLower(email)).Scan(
		&data.ID, &data.Code, &data.Name, &data.Email, &data.Password, &data.RoleID, &data.Status,
		&data.CreatedAt, &data.UpdatedAt, &data.DeletedAt, &data.Role.Code, &data.Role.Name,
	)

	return data, err
}

// Store ...
func (model adminModel) Store(body viewmodel.AdminVM, changedAt time.Time) (res string, err error) {
	sql :=
		`INSERT INTO "admins" (
			"code", "name", "email", "password", "role_id", "status", "created_at", "updated_at"
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $7) RETURNING "id"`
	err = model.DB.QueryRow(sql,
		body.Code, body.Name, body.Email, body.Password, body.RoleID, body.Status, changedAt,
	).Scan(&res)

	return res, err
}

// Update ...
func (model adminModel) Update(id string, body viewmodel.AdminVM, changedAt time.Time) (res string, err error) {
	sql :=
		`UPDATE "admins"
		SET "code" = $1, "name" = $2, "email" = $3, "password" = $4, "role_id" = $5, "status" = $6,
		"updated_at" = $7 WHERE "deleted_at" IS NULL AND "id" = $8 RETURNING "id"`
	err = model.DB.QueryRow(sql,
		body.Code, body.Name, body.Email, body.Password, body.RoleID, body.Status, changedAt, id,
	).Scan(&res)

	return res, err
}

// Destroy ...
func (model adminModel) Destroy(id string, changedAt time.Time) (res string, err error) {
	sql :=
		`UPDATE "admins" SET "updated_at" = $1, "deleted_at" = $1
		WHERE "deleted_at" IS NULL AND "id" = $2 RETURNING "id"`
	err = model.DB.QueryRow(sql, changedAt, id).Scan(&res)

	return res, err
}
