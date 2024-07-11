package transaction

import "Relatdb/meta"

type LogType int

const (
	_ LogType = iota
	TRX_START
	ROLL_BACK
	COMMIT
	ROW
)

type OpType int

const (
	_ OpType = iota
	INSERT
	UPDATE
	DELETE
)

type TrxLog struct {
	lsn       uint64
	logType   LogType
	trxId     uint
	tableName string
	opType    OpType
	before    meta.IndexEntry
	after     meta.IndexEntry
}

func NewTrxLog(
	trxId uint, logType LogType,
	tableName string, opType OpType,
	before meta.IndexEntry, after meta.IndexEntry,
) *TrxLog {
	return &TrxLog{
		lsn:       0,
		logType:   logType,
		trxId:     trxId,
		tableName: tableName,
		opType:    opType,
		before:    before,
		after:     after,
	}
}

func NewTrxStartLog(trxId uint) *TrxLog {
	return NewTrxLog(trxId, TRX_START, "", 0, nil, nil)
}

func NewRowLog(
	trxId uint, tableName string, opType OpType,
	before meta.IndexEntry, after meta.IndexEntry,
) *TrxLog {
	return NewTrxLog(trxId, ROW, tableName, opType, before, after)
}
