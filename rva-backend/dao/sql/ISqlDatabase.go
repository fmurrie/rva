package sql

type ISqlDatabase interface {
	GetDriver() string
	GetHost() string
	GetPort() string
	GetDbName() string
	GetUser() string
	GetPassword() string
	GetCharset() string
	DeployDatabase()
}
