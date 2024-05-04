package meta

type Field struct {
	Index   uint
	Name    string
	Type    byte
	Comment string
	Flag    uint
}

func NewField(index uint, name string, t byte, comment string, flag uint) *Field {
	return &Field{
		Index:   index,
		Name:    name,
		Type:    t,
		Comment: comment,
		Flag:    flag,
	}
}
