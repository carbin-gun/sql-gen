package drivers

import (
	"fmt"
	"strings"

	"github.com/carbin-gun/sql-gen/model"
)

var drivers map[string]Driver

func init() {
	drivers = map[string]Driver{
		"mysql":    MySQLDriver{},
		"postgres": PostgresDriver{},
	}
}

//Driver to define what each driver implementation can do , cover the implementation details
type Driver interface {
	BuildExecutor(host, user, password, database, schema string, ssl bool) Executor
}

//Executor runs to query data
type Executor interface {
	LoadMeta(tables ...string) (map[string]*model.TableMeta, error)
}

//Run Invoke entrace
func Run(host, driver, user, password, database, schema string, ssl bool, tables string) map[string]*model.TableMeta {
	runnerDriver, ok := drivers[driver]
	if !ok {
		panic("driver not support:" + driver)
	}
	e := runnerDriver.BuildExecutor(host, user, password, database, schema, ssl)
	var items []string
	if tables != "" {
		items = strings.Split(tables, ",")
		fmt.Println("query tabels:", items)
	}
	data, err := e.LoadMeta(items...)
	if err != nil {
		msg := fmt.Sprintf("LoadMeta eror:%+v", err)
		panic(msg)
	}
	return data
}
