package corm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Client struct {
	Db *gorm.DB

	Config *Config
}

func (this *Client) IsConnected() bool {
	return this.Db != nil
}

func (this *Client) WithConfig(config *Config) *Client {
	this.Config = config
	return this
}

// get DBOpener for the given driver
func (this *Client) getDBOpener() DBOpener {
	switch this.Config.Driver {
	case DBDriverMySQL:
		return mysql.Open
	case DBDriverSqlite:
		return sqlite.Open
	case DBDriverPostgres:
		return postgres.Open
	default:
		if this.Config.Driver == "" {
			panic("must set corm.config.driver")
		} else {
			panic("unknown database driver: " + this.Config.Driver)
		}

	}
	return nil // unreachable, make compiler happy
}

// now support mysql,sqlite,postgres
func (this *Client) Open() *Client {
	var err error
	if !this.Config.HasDsn() {
		panic("you must set gorm dsn first")
	}
	this.Db, err = gorm.Open(this.getDBOpener()(this.Config.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //使用单数表名，例如 `user` 而不是 `users`
			TablePrefix:   this.Config.TablePrefix,
		},
	})
	if err != nil {
		panic(err.Error())
	}

	//this.Db.Callback().Create().Before("gorm:create").Register("custom_before_create", func(db *gorm.DB) {
	//	db.Statement.Schema.LookUpField("string").DataType = "varchar(255)"
	//})

	//this.Db.Config.NamingStrategy = schema.NamingStrategy{
	//	ColumnMapper: func(column string) string {
	//		if column == "string" {
	//			return "varchar(255)"
	//		}
	//		return column
	//	},
	//}

	if d, err := this.Db.DB(); err != nil {
		panic(err.Error())
	} else {
		d.SetMaxIdleConns(this.Config.MaxOpenConns)
		d.SetMaxOpenConns(this.Config.MaxOpenConns)
		d.SetConnMaxIdleTime(this.Config.ConnMaxIdleTime)
	}
	return this
}

func (this *Client) RegModel(models ...any) error {
	return this.Db.AutoMigrate(models...)
}

func NewClient() (dao *Client) {
	dao = &Client{}
	dao.Config = DefaultConfig()
	return
}

var (
	clients map[string]*Client
)

func ClientIns(name ...string) (client *Client) {
	var (
		ok         bool
		clientName string
	)
	if clients == nil {
		clients = make(map[string]*Client)
	}

	if len(name) > 0 {
		clientName = name[0]
	} else {
		clientName = "default"
	}
	if client, ok = clients[clientName]; !ok {
		client = NewClient()
		clients[clientName] = client
	}
	return
}
