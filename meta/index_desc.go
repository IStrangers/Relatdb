package meta

import "Relatdb/common"

type IndexDesc struct {
	Fields       []*Field
	PrimaryFiled *Field
	FieldMap     map[string]*Field
}

func NewIndexDesc(fields []*Field) *IndexDesc {
	var primaryFiled *Field
	fieldMap := make(map[string]*Field, len(fields))
	for _, field := range fields {
		fieldMap[field.Name] = field
		if field.Flag&common.PRIMARY_KEY_FLAG != 0 {
			primaryFiled = field
		}
	}
	return NewIndexDescByAllArgs(fields, primaryFiled, fieldMap)
}

func NewIndexDescByAllArgs(fields []*Field, primaryFiled *Field, fieldMap map[string]*Field) *IndexDesc {
	desc := &IndexDesc{
		Fields:       fields,
		PrimaryFiled: primaryFiled,
		FieldMap:     fieldMap,
	}
	return desc
}
