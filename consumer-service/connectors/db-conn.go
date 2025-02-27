package connectors

import (
	"consumer-service/logs"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Postgres struct {
	Conn *pgx.Conn
}

var DbConn *Postgres

func GetDatabaseConnection() *Postgres {
	return DbConn
}

func NewPostgres() {
	dsn := "postgres://"
	dsn += viper.GetString("postgresql_user") + ":" + viper.GetString("postgresql_password")
	dsn += "@" + viper.GetString("postgresql_host") + ":" + viper.GetString("postgresql_port")
	dsn += "/" + viper.GetString("postgresql_dbname")

	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		logs.Log.Fatalf("Failed to connect to PostgreSQL: %v\n", err)
	}
	DbConn = &Postgres{Conn: db}
}

// InsertUser inserts a new user record into the users table.
func (p *Postgres) InsertUser(user *User) error {
	query := `
		INSERT INTO users (id,first_name, last_name, email_address, created_at, merged_at, deleted_at, parent_user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := p.Conn.Exec(
		context.Background(),
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.EmailAddress,
		user.CreatedAt,
		user.MergedAt,
		user.DeletedAt,
		user.ParentUserID,
	)
	if err != nil {
		logs.Log.Errorln("Failed to insert user: ", err)
		return err
	}
	return nil
}
