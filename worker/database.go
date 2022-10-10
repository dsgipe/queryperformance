package worker

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"sync"
)

// default string, can be updated with env variables POSTGRES
var connStr = "postgres://postgres:password@localhost:5432/homework"

var conn *pgxpool.Pool

var mutex sync.Mutex

func TimescaleConn() *pgxpool.Pool {
	mutex.Lock()
	if conn == nil {
		if envConnStr, exists := os.LookupEnv("POSTGRES"); exists {
			connStr = envConnStr
		}
		var err error
		conn, err = pgxpool.New(context.Background(), connStr)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
	}
	mutex.Unlock()
	return conn
}
