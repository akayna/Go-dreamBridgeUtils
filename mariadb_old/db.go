package mariadb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

// ConfiguracaoServidor - Estrutura contendo a configuração para acesso ao servidor de DB
type ConfiguracaoServidor struct {
	Usuario string `json:"usuario"`
	Senha   string `json:"senha"`
	URL     string `json:"url"`
	Porta   string `json:"porta"`
	NomeDB  string `json:"dbnome"`
	DB      *sql.DB
}

// Variaveis Globais
var nomeArquivoConfiguracao = "dbconfig"
var dbConf *ConfiguracaoServidor

// AÇÕES GERAIS

// CarregaConfiguracoes - Le os arquivos de configuração de conexão com a DB
func carregaConfiguracoes(nomeArquivo string) (*ConfiguracaoServidor, error) {

	var configuracao ConfiguracaoServidor

	config, err := ioutil.ReadFile(nomeArquivo + ".json")

	if err != nil {
		log.Fatalf("mariadb.carregaConfiguracoes - Erro ao ler arquivo de configuração: %v", err)
		return &configuracao, err
	}

	err = json.Unmarshal(config, &configuracao)

	if err != nil {
		log.Println("mariadb.carregaConfiguracoes - Erro ao ler informações do arquivo de configuração.")
		return &configuracao, err
	}

	//log.Println("mariadb.carregaConfiguracoes - Arquivo de configuração carregado com sucesso: ", configuracao)

	return &configuracao, err
}

func montaURLConexao(configuracao *ConfiguracaoServidor) string {
	//formato da URL => user:password@tcp(127.0.0.1:3306)/database

	if configuracao.URL == "localhost" || configuracao.URL == "127.0.0.1" {
		return configuracao.Usuario + ":" + configuracao.Senha + "@tcp(:" + configuracao.Porta + ")/" + configuracao.NomeDB
	}

	return configuracao.Usuario + ":" + configuracao.Senha + "@tcp(" + configuracao.URL + ":" + configuracao.Porta + ")/" + configuracao.NomeDB

}

// InicializaDB - realiza os processos para inicialização do DB e o mantém aberto.
func InicializaDB() error {
	var err error

	dbConf, err = carregaConfiguracoes(nomeArquivoConfiguracao)

	if err != nil {
		log.Println("mariadb.InicializaDB - Erro ao carregar arquivo de configuração.")
		return err
	}

	dbConf.DB, err = inicializaDB(montaURLConexao(dbConf))

	if err != nil {
		log.Println("mariadb.InicializaDB - Erro ao inicializar banco de dados.")
		return err
	}

	if dbConf.DB == nil {
		log.Panicln("mariadb.InicializaDB - Ponteiro para conexão com DB nulo.")
		return nil
	}
	return nil
}

// Finaliza - Chama todas as finalizações necessárias
func Finaliza() {
	//dbConf = nil
	finalizaDB(dbConf.DB)
}

// ExecutaSelect - Executa uma operação de select no banco de dados
func ExecutaSelect(query string, args ...interface{}) (*sql.Rows, error) {
	// Prepara a query
	stmt, err := dbConf.DB.Prepare(query)

	if err != nil {
		log.Fatalf("mariadb.ExecutaSelect - Erro ao preparar a query: " + query)
		return nil, err
	}

	defer stmt.Close()

	// Executa a query no banco de dados
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("mariadb.ExecutaSelect - Erro ao executar a query: " + query)
		return nil, err
	}

	if rows == nil {
		log.Println("mariadb.ExecutaSelect - rows nulo.")
		return nil, errors.New("mariadb.ExecutaSelect - rows nulo")
	}

	return rows, nil
}

// ExecutaUpdateInsertDelete - Executa uma operação de delete, update ou insert no banco de dados
func ExecutaUpdateInsertDelete(query string, args ...interface{}) (*sql.Result, error) {
	// Prepara a query
	stmt, err := dbConf.DB.Prepare(query)

	if err != nil {
		log.Println("ExecutaQuery - Erro ao preparar a query: " + query)
		log.Printf("ExecutaQuery - Erro: %v", err)
		return nil, err
	}
	defer stmt.Close()

	// Executa a query no banco de dados
	result, err := stmt.Exec(args...)

	if err != nil {
		log.Printf("ExecutaQuery - Erro ao executar a query: " + query)
		fmt.Println(err.Error())
		//log.Fatalf("Erro: %v", err)
		return nil, err
	}

	return &result, nil
}
