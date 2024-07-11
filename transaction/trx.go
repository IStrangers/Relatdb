package transaction

import (
	"Relatdb/meta"
)

type TrxState int

const (
	_ TrxState = iota
	//未开始
	TRX_STATE_NOT_STARTED
	//进行中
	TRX_STATE_ACTIVE
	//2PC/XA
	TRX_STATE_PREPARED
	//已提交
	TRX_STATE_COMPLETED
)

type Trx struct {
	trxId    uint
	state    TrxState
	logs     []*TrxLog
	logStore *LogStore
}

func (self *Trx) AddLogByTrxLog(log *TrxLog) {
	self.logs = append(self.logs, log)
}

func checkIndexEntry(indexEntry meta.IndexEntry) {
	if indexEntry != nil {
		_, ok := indexEntry.(*meta.ClusterIndexEntry)
		if !ok {
			panic("log before must be of type ClusterIndexEntry")
		}
	}
}

func (self *Trx) AddLog(tableName string, opType OpType, before meta.IndexEntry, after meta.IndexEntry) {
	checkIndexEntry(before)
	checkIndexEntry(after)
	log := NewRowLog(self.trxId, tableName, opType, before, after)
	self.AddLogByTrxLog(log)
}

func (self *Trx) Redo() {

}

func (self *Trx) Undo() {

}

func (self *Trx) Begin() {
	log := NewTrxStartLog(self.trxId)
	self.logStore.AppendLog(log)
	self.state = TRX_STATE_ACTIVE
}

func (self *Trx) Rollback() {

}

func (self *Trx) Commit() {

}
