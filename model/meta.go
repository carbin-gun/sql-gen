package model

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/stoewer/go-strcase"
)

//TableMeta table meta data
type TableMeta struct {
	TableName string
	Database  string
	Schema    string
	Columns   []*ColumnMeta
}

//JointColumns returns the xx,xx,xx format string
func (t *TableMeta) JointColumns() string {
	buf := bytes.Buffer{}
	size := len(t.Columns)
	for i, c := range t.Columns {
		buf.WriteString(string(c.Column))
		if i != size-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

//StructName returns the struct name according to a rule of
// aaa_bbb --> AaaBbb
func (t *TableMeta) StructName() string {
	items := strings.Split(t.TableName, "_")
	buf := bytes.Buffer{}
	for _, item := range items {
		buf.WriteString(strings.Title(item)) //Title it
	}
	return buf.String()
}

type ColumnName string

func (c ColumnName) CamelCase() string {
	return strcase.UpperCamelCase(string(c))
}

type ColumnType string

const (
	BigInt            ColumnType = "bigint"
	Int               ColumnType = "int"
	SmallInt          ColumnType = "smallint"
	Varchar           ColumnType = "character varying"
	JSONB             ColumnType = "jsonb"
	TimestampWithZone ColumnType = "timestamp with time zone"
	Boolen            ColumnType = "boolean"
	Integer           ColumnType = "integer"
	Date              ColumnType = "date"
)

var (
	SqlTypeToGoType = map[ColumnType]string{
		BigInt:            "int64",
		Int:               "int64",
		SmallInt:          "uint8",
		Varchar:           "string",
		JSONB:             "interface{}",
		TimestampWithZone: "time.Time",
		Boolen:            "bool",
		Integer:           "int64",
		Date:              "time.Time",
	}
)

func (t ColumnType) GoType() string {
	target, exist := SqlTypeToGoType[t]
	if exist {
		return target
	}
	return string(t)
}

//ColumnMeta ,column meta data
type ColumnMeta struct {
	TableName string     `db:"table_name"`
	Database  string     `db:"table_catalog"`
	Schema    string     `db:"table_schema"`
	Column    ColumnName `db:"column_name"`
	Type      ColumnType `db:"data_type"`
	Ordinal   int        `db:"ordinal_position"`
}

func (c *ColumnMeta) String() string {
	return fmt.Sprintf("table:%s,column:%s", c.TableName, c.Column)
}
