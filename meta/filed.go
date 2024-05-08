package meta

type Field struct {
	Index        uint
	Name         string
	Type         byte
	Flag         uint
	DefaultValue Value
	Comment      string
}

func NewField(index uint, name string, t byte, flag uint, defaultValue Value, comment string) *Field {
	return &Field{
		Index:        index,
		Name:         name,
		Type:         t,
		Flag:         flag,
		DefaultValue: defaultValue,
		Comment:      comment,
	}
}
