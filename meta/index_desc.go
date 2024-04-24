package meta

type Attribute struct {
	Index        uint
	Name         string
	Type         byte
	Comment      string
	IsPrimaryKey bool
}

func CreateAttribute(index uint, name string, t byte, comment string) *Attribute {
	return &Attribute{
		Index:   index,
		Name:    name,
		Type:    t,
		Comment: comment,
	}
}

type IndexDesc struct {
	Attrs       []*Attribute
	PrimaryAttr *Attribute
	AttrMap     map[string]*Attribute
}

func CreateIndexDesc(attrs []*Attribute) *IndexDesc {
	desc := &IndexDesc{
		Attrs: make([]*Attribute, 0),
	}
	desc.AttrMap = make(map[string]*Attribute, len(attrs))
	for _, attr := range attrs {
		desc.AttrMap[attr.Name] = attr
		if attr.IsPrimaryKey {
			desc.PrimaryAttr = attr
		}
	}
	return desc
}
