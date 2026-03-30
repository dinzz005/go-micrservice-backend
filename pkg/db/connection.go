package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func NewDB() *sql.DB {
	 //load env
	 err := godotenv.Load()
	 if err != nil {
		 log.Fatal("Error Loading config file")
	 }

	 cfg := mysql.Config{
		 User: os.Getenv("DB_USER"),
		 Passwd: os.Getenv("DB_PASSWORD"),
		 Net: "tcp",
		 Addr: os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		 DBName: os.Getenv("DB_NAME"),
		 ParseTime: true,
		 AllowNativePasswords: true,
		 Loc: time.Local,
	 }
	 dsn := cfg.FormatDSN()

	 db, err := sql.Open("mysql", dsn)
	 if err != nil {
		 log.Fatal(err)
	 }

	 if err := db.Ping(); err != nil {
		 log.Fatal(err)
	 }
   
	 return  db
}
