package meta

type IndexDesc struct {
	Fields       []*Field
	PrimaryFiled *Field
	FieldMap     map[string]*Field
}

func NewIndexDesc(fields []*Field) *IndexDesc {
	desc := &IndexDesc{
		Fields: make([]*Field, 0),
	}
	desc.FieldMap = make(map[string]*Field, len(fields))
	for _, field := range fields {
		desc.FieldMap[field.Name] = field
		if field.IsPrimaryKey {
			desc.PrimaryFiled = field
		}
	}
	return desc
}
