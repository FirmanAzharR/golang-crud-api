package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Konfigurasi koneksi PostgreSQL
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "go_db_crud"
)

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100) UNIQUE
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Gagal membuat tabel:", err)
	}

	fmt.Println("Tabel users berhasil dibuat!")
}

func createUser(db *sql.DB, name, email string) {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	var id int
	err := db.QueryRow(query, name, email).Scan(&id)
	if err != nil {
		log.Fatal("Gagal menambahkan user:", err)
	}
	fmt.Println("User berhasil ditambahkan dengan ID:", id)
}

func getUsers(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Fatal("Gagal mengambil data:", err)
	}
	defer rows.Close()

	fmt.Println("Daftar Users:")
	for rows.Next() {
		var id int
		var name, email string
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", id, name, email)
	}
}

func getUserByID(db *sql.DB, userID int) {
	var name, email string
	query := `SELECT name, email FROM users WHERE id=$1`
	err := db.QueryRow(query, userID).Scan(&name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User tidak ditemukan")
		} else {
			log.Fatal(err)
		}
		return
	}
	fmt.Printf("User ditemukan: %s - %s\n", name, email)
}

func updateUser(db *sql.DB, id int, name, email string) {
	query := `UPDATE users SET name=$1, email=$2 WHERE id=$3`
	_, err := db.Exec(query, name, email, id)
	if err != nil {
		log.Fatal("Gagal memperbarui user:", err)
	}
	fmt.Println("User berhasil diperbarui!")
}

// func deleteUser(db *sql.DB, id int) {
// 	query := `DELETE FROM users WHERE id=$1`
// 	_, err := db.Exec(query, id)
// 	if err != nil {
// 		log.Fatal("Gagal menghapus user:", err)
// 	}
// 	fmt.Println("User berhasil dihapus!")
// }

func main() {
	// Format koneksi PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// fmt.Println(psqlInfo)
	// Membuka koneksi
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cek koneksi
	err = db.Ping()
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	fmt.Println("Berhasil terhubung ke database!")

	// createTable(db)                                 // Membuat tabel jika belum ada
	// createUser(db, "firman", "firman@example.com") // Menambahkan data
	getUsers(db) // Menampilkan semua data
	// getUserByID(db, 1)                              // Menampilkan user tertentu
	// updateUser(db, 1, "Bobby", "bobby@example.com") // Update data
	// deleteUser(db, 1)                               // Hapus data

}
