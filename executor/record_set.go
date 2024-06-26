package executor

import "Relatdb/meta"

type RecordSet interface {
	GetAffectedRows() uint64
	GetInsertId() uint64
	GetColumns() []meta.Value
	GetRows() [][]meta.Value
}

type RecordSetImpl struct {
	affectedRows uint64
	insertId     uint64
	columns      []meta.Value
	rows         [][]meta.Value
}

func NewRecordSet(affectedRows uint64, insertId uint64, columns []meta.Value, rows [][]meta.Value) RecordSet {
	return &RecordSetImpl{
		affectedRows: affectedRows,
		insertId:     insertId,
		columns:      columns,
		rows:         rows,
	}
}

func (self *RecordSetImpl) GetAffectedRows() uint64 {
	return self.affectedRows
}

func (self *RecordSetImpl) GetInsertId() uint64 {
	return self.insertId
}

func (self *RecordSetImpl) GetColumns() []meta.Value {
	return self.columns
}

func (self *RecordSetImpl) GetRows() [][]meta.Value {
	return self.rows
}
