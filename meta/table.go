package meta

import (
	"Relatdb/common"
	"Relatdb/index"
)

type Table struct {
	MetaPath         string
	DataPath         string
	Name             string
	Fields           []*Field
	PrimaryFiled     *Field
	FieldMap         map[string]uint
	ClusterIndex     *index.BaseIndex
	SecondaryIndexes []*index.BaseIndex
}

func NewTable(name string, fields []*Field) *Table {
	fieldMap := make(map[string]uint, len(fields))
	var primaryFiled *Field
	for i, field := range fields {
		if field.Flag&common.PRIMARY_KEY_FLAG != 0 {
			primaryFiled = field
		}
		fieldMap[field.Name] = uint(i)
	}
	return &Table{
		Name:         name,
		Fields:       fields,
		PrimaryFiled: primaryFiled,
		FieldMap:     fieldMap,
	}
}
