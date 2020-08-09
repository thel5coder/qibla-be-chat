package model

import (
	"qibla-backend-chat/usecase/viewmodel"
	"time"

	"database/sql"
)

const (
	// FileGroupRoomProfilePicture ...
	FileGroupRoomProfilePicture = "room_profile_picture"
)

// IFile ...
type IFile interface {
	FindAllByUserID(userID, types string) (data []FileEntity, err error)
	FindByID(id string) (FileEntity, error)
	FindUnassignedByID(id, types, userID string) (FileEntity, error)
	FindAssignedByID(id, types string) (FileEntity, error)
	Store(body viewmodel.FileVM, changedAt time.Time) (string, error)
	Destroy(id string, changedAt time.Time) (string, error)
}

// FileEntity ....
type FileEntity struct {
	ID        string         `db:"id"`
	Type      sql.NullString `db:"type"`
	Path      sql.NullString `db:"path"`
	UserID    sql.NullString `db:"user_id"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// fileModel ...
type fileModel struct {
	DB *sql.DB
}

// NewFileModel ...
func NewFileModel(db *sql.DB) IFile {
	return &fileModel{DB: db}
}

// FindAllByUserID ...
func (model fileModel) FindAllByUserID(userID, types string) (data []FileEntity, err error) {
	query := `SELECT "id", "type", "path", "user_id", "created_at", "updated_at", "deleted_at"
		FROM "files" WHERE "deleted_at" IS NULL AND "user_id" = $1 AND "type" = $2
		ORDER BY "created_at"`

	rows, err := model.DB.Query(query, userID, types)
	if err != nil {
		return data, err
	}

	defer rows.Close()
	for rows.Next() {
		d := FileEntity{}
		err = rows.Scan(
			&d.ID, &d.Type, &d.Path, &d.UserID, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
		)
		if err != nil {
			return data, err
		}
		data = append(data, d)
	}
	err = rows.Err()

	return data, err
}

// FindByID ...
func (model fileModel) FindByID(id string) (data FileEntity, err error) {
	query :=
		`SELECT "id", "type", "path", "user_id", "created_at", "updated_at", "deleted_at"
		FROM "files" WHERE "deleted_at" IS NULL AND "id" = $1
		ORDER BY "created_at" DESC LIMIT 1`
	err = model.DB.QueryRow(query, id).Scan(
		&data.ID, &data.Type, &data.Path, &data.UserID, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
	)

	return data, err
}

// FindUnassignedByID ...
func (model fileModel) FindUnassignedByID(id, types, userID string) (data FileEntity, err error) {
	query :=
		`SELECT "id", "type", "path", "user_id", "created_at", "updated_at", "deleted_at"
		FROM "files" WHERE "deleted_at" IS NULL AND "id" = $1 AND "type" = $2
		ORDER BY "created_at" DESC LIMIT 1`
	err = model.DB.QueryRow(query, id, types).Scan(
		&data.ID, &data.Type, &data.Path, &data.UserID, &data.CreatedAt, &data.UpdatedAt,
		&data.DeletedAt,
	)

	return data, err
}

// FindAssignedByID ...
func (model fileModel) FindAssignedByID(id, types string) (data FileEntity, err error) {
	query :=
		`SELECT uf."id", uf."type", uf."path", uf."user_id", uf."created_at", uf."updated_at", uf."deleted_at"
		FROM "files" uf
		WHERE uf."deleted_at" IS NULL AND uf."id" = $1 AND uf."type" = $2
		ORDER BY uf."created_at" DESC LIMIT 1`
	err = model.DB.QueryRow(query, id, types).Scan(
		&data.ID, &data.Type, &data.Path, &data.UserID, &data.CreatedAt, &data.UpdatedAt,
		&data.DeletedAt,
	)

	return data, err
}

// Store ...
func (model fileModel) Store(body viewmodel.FileVM, changedAt time.Time) (res string, err error) {
	sql :=
		`INSERT INTO "files" ("type", "path", "user_id", "created_at", "updated_at")
		VALUES($1, $2, $3, $4, $4) RETURNING "id"`
	err = model.DB.QueryRow(sql, body.Type, body.Path, body.UserID, changedAt).Scan(&res)

	return res, err
}

// Destroy ...
func (model fileModel) Destroy(id string, changedAt time.Time) (res string, err error) {
	sql := `UPDATE "files" SET deleted_at = $1 WHERE id = $2 RETURNING "id"`
	err = model.DB.QueryRow(sql, changedAt, id).Scan(&res)

	return res, err
}
