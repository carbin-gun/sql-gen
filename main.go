package main

import (
	"flag"
	"fmt"
	u "os/user"

	"github.com/carbin-gun/sql-gen/drivers"
	"github.com/carbin-gun/sql-gen/writers"
	"github.com/pkg/errors"
)

var (
	host     = flag.String("host", "localhost", "is the database postgres")
	ssl      = flag.Bool("ssl", false, "database access ssl on/off")
	driver   = flag.String("driver", "", "the driver<mysql/postgres> for access")
	user     = flag.String("user", "", "user for the databse")
	password = flag.String("password", "", "password for the databse")
	database = flag.String("database", "", "database for the access")
	schema   = flag.String("schema", "public", "database for the access")
	tables   = flag.String("tables", "", "tables for the access,default for all,split by `,` is spcefied ")
	writer   = flag.String("writer", "file", "output to file/console")
)

func main() {
	flag.Parse()
	if *database == "" {
		panic("you must specify a database ")
	}
	if *driver == "" {
		panic("you must specify a driver,mysql or postgres")
	}
	w, ok := writers.Lookup(*writer)
	if !ok {
		panic("you can only specify writer to be file/console")
	}
	if *user == "" {
		current, err := u.Current()
		fmt.Printf("user current user <%s> for db<%s> login. \n", current.Username, *database)
		if err != nil {
			err = errors.Wrap(err, "os/user.Current() error")
			panic(err)
		}
		*user = current.Username
	}
	fmt.Println("tables:",*tables)
	data := drivers.Run(*host, *driver, *user, *password, *database, *schema, *ssl, *tables)
	err := w.Write(data)
	if err != nil {
		panic(err)
	}
}
