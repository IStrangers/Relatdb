package meta

type DataBase struct {
	Name     string
	TableMap map[string]*Table
}

func NewDataBase(name string) *DataBase {
	return &DataBase{
		Name:     name,
		TableMap: make(map[string]*Table),
	}
}

func (self *DataBase) GetTable(tableName string) *Table {
	table := self.TableMap[tableName]
	return table
}
