package meta

type Table struct {
	MetaPath         string
	DataPath         string
	Name             string
	Fields           []*Field
	PrimaryFiled     *Field
	FieldMap         map[string]uint
	ClusterIndex     Index
	SecondaryIndexes []Index
}

func NewTable(
	name string, fields []*Field, primaryFiled *Field,
	fieldMap map[string]uint, clusterIndex Index,
	secondaryIndexes []Index,
) *Table {
	return &Table{
		Name:             name,
		Fields:           fields,
		PrimaryFiled:     primaryFiled,
		FieldMap:         fieldMap,
		ClusterIndex:     clusterIndex,
		SecondaryIndexes: secondaryIndexes,
	}
}
