package drivers

import "github.com/carbin-gun/sql-gen/model"

func transform(columns []*model.ColumnMeta) map[string]*model.TableMeta {
	tables := map[string]*model.TableMeta{}
	for i := range columns {
		column := columns[i]
		meta, ok := tables[column.TableName]
		if !ok {
			tables[column.TableName] = &model.TableMeta{
				TableName: column.TableName,
				Database:  column.Database,
				Schema:    column.Schema,
				Columns:   []*model.ColumnMeta{column},
			}
			continue
		}
		//already exists
		meta.Columns = append(meta.Columns, column)
	}
	return tables
}
