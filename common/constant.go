package common

/*
 * FieldType
 */
const (
	FIELD_TYPE_DECIMAL     = 0
	FIELD_TYPE_TINY        = 1
	FIELD_TYPE_SHORT       = 2
	FIELD_TYPE_LONG        = 3
	FIELD_TYPE_FLOAT       = 4
	FIELD_TYPE_DOUBLE      = 5
	FIELD_TYPE_NULL        = 6
	FIELD_TYPE_TIMESTAMP   = 7
	FIELD_TYPE_LONGLONG    = 8
	FIELD_TYPE_INT24       = 9
	FIELD_TYPE_DATE        = 10
	FIELD_TYPE_TIME        = 11
	FIELD_TYPE_DATETIME    = 12
	FIELD_TYPE_YEAR        = 13
	FIELD_TYPE_NEWDATE     = 14
	FIELD_TYPE_VARCHAR     = 15
	FIELD_TYPE_BIT         = 16
	FIELD_TYPE_NEW_DECIMAL = 246
	FIELD_TYPE_ENUM        = 247
	FIELD_TYPE_SET         = 248
	FIELD_TYPE_TINY_BLOB   = 249
	FIELD_TYPE_MEDIUM_BLOB = 250
	FIELD_TYPE_LONG_BLOB   = 251
	FIELD_TYPE_BLOB        = 252
	FIELD_TYPE_VAR_STRING  = 253
	FIELD_TYPE_STRING      = 254
	FIELD_TYPE_GEOMETRY    = 255
)

/*
 * FieldFlag
 */
const (
	NOT_NULL_FLAG       = 0x0001
	PRI_KEY_FLAG        = 0x0002
	UNIQUE_KEY_FLAG     = 0x0004
	MULTIPLE_KEY_FLAG   = 0x0008
	BLOB_FLAG           = 0x0010
	UNSIGNED_FLAG       = 0x0020
	ZEROFILL_FLAG       = 0x0040
	BINARY_FLAG         = 0x0080
	ENUM_FLAG           = 0x0100
	AUTO_INCREMENT_FLAG = 0x0200
	TIMESTAMP_FLAG      = 0x0400
	SET_FLAG            = 0x0800
)

type LengthAndDecimal struct {
	Length  int
	Decimal int
}

var defaultLengthAndDecimal = map[byte]LengthAndDecimal{
	FIELD_TYPE_BIT:         {1, 0},
	FIELD_TYPE_TINY:        {4, 0},
	FIELD_TYPE_SHORT:       {6, 0},
	FIELD_TYPE_INT24:       {9, 0},
	FIELD_TYPE_LONG:        {11, 0},
	FIELD_TYPE_LONGLONG:    {20, 0},
	FIELD_TYPE_DOUBLE:      {22, -1},
	FIELD_TYPE_FLOAT:       {12, -1},
	FIELD_TYPE_NEW_DECIMAL: {10, 0},
	FIELD_TYPE_DATE:        {10, 0},
	FIELD_TYPE_TIMESTAMP:   {19, 0},
	FIELD_TYPE_DATETIME:    {19, 0},
	FIELD_TYPE_STRING:      {1, 0},
	FIELD_TYPE_VARCHAR:     {5, 0},
	FIELD_TYPE_VAR_STRING:  {5, 0},
	FIELD_TYPE_TINY_BLOB:   {255, 0},
	FIELD_TYPE_BLOB:        {65535, 0},
	FIELD_TYPE_MEDIUM_BLOB: {16777215, 0},
	FIELD_TYPE_LONG_BLOB:   {4294967295, 0},
	FIELD_TYPE_NULL:        {0, 0},
	FIELD_TYPE_SET:         {-1, 0},
	FIELD_TYPE_ENUM:        {-1, 0},
}

func GetFieldDefaultLengthAndDecimal(fieldType byte) LengthAndDecimal {
	return defaultLengthAndDecimal[fieldType]
}
