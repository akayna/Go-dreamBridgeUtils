package mariadb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" //Necessário para abertura da conexão com o DB
)

// MaxConnections - Número máximo de conexões abertas no DB
var MaxConnections = 3

// IdleTimeMinutes - Tempo em minutos para manter conexões aberta no estado IDLE
const IdleTimeMinutes int = 10

// InicializaDB - Realiza uma conexão com a base de dados e retorna a conexão.
func inicializaDB(urlConexao string) (*sql.DB, error) {
	db, err := sql.Open("mysql", urlConexao)

	if err != nil {
		log.Fatalf("mariadb.inicializaDB - Erro ao abrir conexao com DB: %v", err)
		return nil, err
	}

	db.SetMaxIdleConns(MaxConnections) // Número máximo de conexões abertas sem uso
	db.SetMaxOpenConns(MaxConnections) // Número máximo de conexões abertas simultaneamente

	var tempoDuracao = time.Minute * time.Duration(IdleTimeMinutes)
	db.SetConnMaxLifetime(tempoDuracao) // Tempo máximo que uma conexão pode ser reutilizada. 0 = para sempre

	err = testaConexaoDB(db)

	if err != nil {
		db = nil
		log.Printf("mariadb.inicializaDB - Falha ao testar conexao com DB.")
		return nil, err
	}

	return db, nil
}

// TestaConexaoDB - Envia um PING para o DB verificar se a conexão está ok.
func testaConexaoDB(db *sql.DB) error {

	if db == nil {
		return errors.New("mariadb.testaConexaoDB - Ponteiro para conexão com o DB nulo")
	}

	err := db.Ping()

	if err != nil {
		log.Println("mariadb.testaConexaoDB - Erro testar conexao com DB.")
		return err
	}

	return nil
}

// ListaTabelasDB - Testa a conexão com o banco de dados.
func listaTabelasDB(db *sql.DB) error {

	if db == nil {
		return errors.New("mariadb.listaTabelasDB - Ponteiro para conexão com o DB nulo")
	}

	rows, err := db.Query("SHOW TABLES")

	if err != nil {
		log.Fatalf("Erro: %v", err)
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Println(name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func finalizaDB(db *sql.DB) error {
	if db == nil {
		return errors.New("mariadb.finalizaDB - Ponteiro para conexão com o DB nulo")
	}
	err := db.Close()
	db = nil

	if err != nil {
		log.Println("mariadb.finalizaDB - Erro ao finalizar banco da dados.")
		return err
	}

	return nil
}
