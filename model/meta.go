package model

import (
	"bytes"
	"fmt"
	"strings"
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
		buf.WriteString(c.Column)
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
		buf.WriteString(strings.Title(item))
	}
	return buf.String()
}

//ColumnMeta ,column meta data
type ColumnMeta struct {
	TableName string `db:"table_name"`
	Database  string `db:"table_catalog"`
	Schema    string `db:"table_schema"`
	Column    string `db:"column_name"`
	Type      string `db:"data_type"`
	Ordinal   int    `db:"ordinal_position"`
}

func (c *ColumnMeta) String() string {
	return fmt.Sprintf("table:%s,column:%s", c.TableName, c.Column)
}
