package meta

type Field struct {
	Index        uint
	Name         string
	Type         byte
	Comment      string
	IsPrimaryKey bool
}

func NewField(index uint, name string, t byte, comment string) *Field {
	return &Field{
		Index:   index,
		Name:    name,
		Type:    t,
		Comment: comment,
	}
}
