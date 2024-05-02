package meta

import "Relatdb/index"

type Table struct {
	MetaPath         string
	TablePath        string
	Name             string
	Fields           []*Field
	PrimaryFiled     *Field
	FieldMap         map[string]uint
	ClusterIndex     *index.BaseIndex
	SecondaryIndexes []*index.BaseIndex
}
