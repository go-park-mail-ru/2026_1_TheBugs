package dsn

import "fmt"

func BuildDSN(host string, port int, user string, password string, dbname string, sslmode string) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslmode,
	)
}
