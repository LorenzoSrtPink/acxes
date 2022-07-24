package acxes

import (
	"os"

	"github.com/jmoiron/sqlx"
)

type serverOption int

const (
	localhost serverOption = iota
	remoteInternalHost
	remoteHost
	customHost
)

type credential struct {
	local          string
	remote         string
	outside        string
	dataBaseEngine string
	database       string
	port           string
	charset        string
	user           string
	password       string
	accessString   string
	host           string
}

var ip string

func (data *credential) iniTokens() {
	data.local = "localhost"
	data.remote = os.Getenv("HostIP")
	data.outside = os.Getenv("RemoteIP")
	data.dataBaseEngine = "mysql"
	data.database = os.Getenv("MySQLDatabase")
	data.port = "3306"
	data.charset = "utf8"
	data.user = os.Getenv("MySQLUser")
	data.password = os.Getenv("MySQLPassword")
	data.host = "none"
}

//AccessLocal for accesing local database
func Local() (*sqlx.DB, error) {
	return access(localhost)
}

//AccessLocalFor temp access for local network
func LocalFor(database string) (*sqlx.DB, error) {
	return accessForDatabase("localhost", database)
}

//AccessRemote for local network server
func Remote() (*sqlx.DB, error) {
	return access(remoteInternalHost)
}

//AccessRemoteOutside for static ISP provider IP adress
func RemoteOutside() (*sqlx.DB, error) {
	return access(remoteHost)
}

//AccessFor any ip
func For(host string) (*sqlx.DB, error) {
	ip = host
	return access(customHost)
}

func access(host serverOption) (*sqlx.DB, error) {
	token := credential{}
	token.iniTokens()

	switch host {
	case localhost:
		token.host = token.local
	case remoteInternalHost:
		token.host = token.remote
	case remoteHost:
		token.host = token.outside
	case customHost:
		token.host = ip

	}

	accessString := token.user + ":" + token.password + "@(" + token.host + ":" + token.port + ")/" + token.database + "?charset=" + token.charset
	return sqlx.Connect(token.dataBaseEngine, accessString)
}

func accessForDatabase(host string, database string) (*sqlx.DB, error) {
	token := credential{}
	token.iniTokens()
	token.host = host
	accessString := token.user + ":" + token.password + "@(" + token.host + ":" + token.port + ")/" + database + "?charset=" + token.charset
	return sqlx.Connect(token.dataBaseEngine, accessString)
}
