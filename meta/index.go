package meta

type Index interface {
	IsPrimary() bool
	IsUnique() bool
}
