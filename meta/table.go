package meta

type Table struct {
	MetaPath         string
	DataPath         string
	DatabaseName     string
	Name             string
	Fields           []*Field
	PrimaryFiled     *Field
	FieldMap         map[string]uint
	ClusterIndex     Index
	SecondaryIndexes []Index
}

func NewTable(
	databaseName string, name string, fields []*Field, primaryFiled *Field,
	fieldMap map[string]uint, clusterIndex Index,
	secondaryIndexes []Index,
) *Table {
	return &Table{
		DatabaseName:     databaseName,
		Name:             name,
		Fields:           fields,
		PrimaryFiled:     primaryFiled,
		FieldMap:         fieldMap,
		ClusterIndex:     clusterIndex,
		SecondaryIndexes: secondaryIndexes,
	}
}
