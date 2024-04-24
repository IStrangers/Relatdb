package meta

type Value interface {
	GetLength() uint
	Compare(Value) int8
}
