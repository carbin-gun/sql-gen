package drivers

import (
	"fmt"

	"github.com/carbin-gun/sql-gen/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres driver
	"github.com/pkg/errors"
)

var (
	//all tables
	metaSQL = `select table_catalog,table_schema, table_name,column_name,data_type,ordinal_position from information_schema.columns 
	where  table_catalog = '%s' and table_schema='%s'   order by table_name,ordinal_position asc`
	//query only specified tables
	metaTablesSQL = `select table_catalog,table_schema, table_name,column_name,data_type,ordinal_position from information_schema.columns 
	where  table_catalog = '%s' and table_schema='%s' and table_name in (?)  order by table_name,ordinal_position asc`
)

//PostgresDriver implements Driver interface
type PostgresDriver struct{}

//PostgresExecutor for Executor interface
type PostgresExecutor struct {
	Host     string
	User     string
	Password string
	DBName   string
	Schema   string
	SSL      string
}

//BuildExecutor Driver interface implementation for postgres
func (p PostgresDriver) BuildExecutor(host, user, password, database, schema string, ssl bool) Executor {
	e := new(PostgresExecutor)
	e.Host = host
	e.User = user
	e.Password = password
	e.DBName = database
	e.Schema = schema
	if ssl {
		e.SSL = "enable"
	} else {
		e.SSL = "disable"
	}
	return e
}

//LoadMeta ...
func (e *PostgresExecutor) LoadMeta(tables ...string) (map[string]*model.TableMeta, error) {
	dsn := fmt.Sprintf("host='%s' user='%s' password='%s' dbname='%s' sslmode='%s'", e.Host, e.User, e.Password, e.DBName, e.SSL)
	fmt.Println("dsn:", dsn)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		err = errors.Wrap(err, "open connection error")
		panic(err)
	}
	defer db.Close()
	result := []*model.ColumnMeta{}
	if len(tables) == 0 { //query all
		sql := fmt.Sprintf(metaSQL, e.DBName, e.Schema)
		err = db.Select(&result, sql)
	} else { //query specified tables
		sql := fmt.Sprintf(metaTablesSQL, e.DBName, e.Schema)
		query, args, err := sqlx.In(sql, tables)
		if err != nil {
			return nil, errors.Wrap(err, "query specified tables error")
		}
		query = db.Rebind(query)
		err = db.Select(&result, query, args...)
	}
	if err != nil {
		return nil, errors.Wrap(err, "postgres LoadMeta error")
	}
	return transform(result), nil
}
