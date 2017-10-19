package drivers

import "github.com/carbin-gun/sql-gen/model"

//MySQLDriver for mysql
type MySQLDriver struct{}

//MySQLExecutor implements executor interface
type MySQLExecutor struct {
}

//BuildExecutor implements driver interface
func (MySQLDriver) BuildExecutor(host, user, password, database, schema string, ssl bool) Executor {
	e := new(MySQLExecutor)
	return e
}

//LoadMeta implement Executor interface
func (e *MySQLExecutor) LoadMeta(tables ...string) (map[string]*model.TableMeta, error) {
	return nil, nil
}
