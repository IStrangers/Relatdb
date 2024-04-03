package parser

import "Relatdb/common"

type Token int

func (tkn Token) String() string {
	return tokenStringMap[tkn]
}

const (
	_ Token = iota
	ILLEGAL
	EOF
	COMMENT
	MULTI_COMMENT
	WHITE_SPACE
	IDENTIFIER

	LEFT_PARENTHESIS  // (
	RIGHT_PARENTHESIS // )
	LEFT_BRACE        // {
	RIGHT_BRACE       // }
	LEFT_BRACKET      // [
	RIGHT_BRACKET     // ]
	DOT               // .
	COMMA             // ,
	COLON             // :
	SEMICOLON         // ;

	NUMBER
	STRING
	BOOLEAN
	NULL

	ADDITION              // +
	SUBTRACT              // -
	MULTIPLY              // *
	DIVIDE                // /
	REMAINDER             // %
	AND_ARITHMETIC        // &
	OR_ARITHMETIC         // |
	INCREMENT             // ++
	DECREMENT             // --
	ADDITION_ASSIGN       // +=
	SUBTRACT_ASSIGN       // -=
	MULTIPLY_ASSIGN       // *=
	DIVIDE_ASSIGN         // /=
	REMAINDER_ASSIGN      // %=
	AND_ARITHMETIC_ASSIGN // &=
	OR_ARITHMETIC_ASSIGN  // |=

	ASSIGN           // =
	EQUAL            // ==
	NOT_ARITHMETIC   // ！
	NOT_EQUAL        // !=
	LESS             // <
	LESS_OR_EQUAL    // <=
	GREATER          // >
	GREATER_OR_EQUAL // >=
	LOGICAL_AND      // &&
	LOGICAL_OR       // ||

	SHOW           // show
	DATABASES      // databases
	TABLES         // tables
	VARIABLES      // variables
	IF             // if
	NOT            // not
	EXISTS         // exists
	CREATE         // create
	DATABASE       // database
	TABLE          // table
	INDEX          // index
	UNIQUE         // unique
	SPATIAL        // spatial
	FULLTEXT       // fulltext
	DROP           // drop
	INTO           // into
	VALUES         // values
	INSERT         // insert
	DELETE         // delete
	UPDATE         // update
	SELECT         // select
	FROM           // from
	WHERE          // where
	GROUP          // group
	HAVING         // having
	ORDER          // order
	LIMIT          // limit
	AS             // as
	SET            // set
	BY             // by
	AES            // aes
	DESC           // desc
	INNER          // inner
	LEFT           // left
	RIGHT          // right
	JOIN           // join
	ON             // on
	IN             // in
	AND            // and
	OR             // or
	BETWEEN        // between
	LIKE           // like
	DISTINCT       // distinct
	UNION          // union
	ALL            // all
	TRUNCATE       // truncate
	PRIMARY        // primary
	KEY            // key
	UNSIGNED       // unsigned
	ZEROFILL       // zerofill
	AUTO_INCREMENT // auto_increment
	DEFAULT        // default
	COLUMN_COMMENT // comment

	TINYINT   // tinyint
	SMALLINT  // smallint
	MEDIUMINT // mediumint
	INT       // int
	INTEGER   // integer
	BIGINT    // bigint
	FLOAT     // float
	DOUBLE    // double
	DECIMAL   // decimal
	DATE      // date
	TIME      // time
	DATETIME  // datetime
	TIMESTAMP // timestamp
	CHAR      // char
	VARCHAR   // varchar
	BINARY    // binary
	VARBINARY // varbinary
	TEXT      // text
	BLOB      // blob
)

var tokenStringMap = [...]string{
	0:             "UNKNOWN",
	ILLEGAL:       "ILLEGAL",
	EOF:           "EOF",
	COMMENT:       "COMMENT",
	MULTI_COMMENT: "MULTI_COMMENT",
	WHITE_SPACE:   "WHITE_SPACE",
	IDENTIFIER:    "IDENTIFIER",

	LEFT_PARENTHESIS:  "(",
	RIGHT_PARENTHESIS: ")",
	LEFT_BRACE:        "{",
	RIGHT_BRACE:       "}",
	LEFT_BRACKET:      "[",
	RIGHT_BRACKET:     "]",
	DOT:               ".",
	COMMA:             ",",
	COLON:             ":",
	SEMICOLON:         ";",

	NUMBER:  "NUMBER",
	STRING:  "STRING",
	BOOLEAN: "BOOLEAN",
	NULL:    "NULL",

	ADDITION:              "+",
	SUBTRACT:              "-",
	MULTIPLY:              "*",
	DIVIDE:                "/",
	REMAINDER:             "%",
	AND_ARITHMETIC:        "&",
	OR_ARITHMETIC:         "|",
	INCREMENT:             "++",
	DECREMENT:             "--",
	ADDITION_ASSIGN:       "+=",
	SUBTRACT_ASSIGN:       "-=",
	MULTIPLY_ASSIGN:       "*=",
	DIVIDE_ASSIGN:         "/=",
	REMAINDER_ASSIGN:      "%=",
	AND_ARITHMETIC_ASSIGN: "&=",
	OR_ARITHMETIC_ASSIGN:  "|=",

	ASSIGN:           "=",
	EQUAL:            "==",
	NOT_ARITHMETIC:   "!",
	NOT_EQUAL:        "!=",
	LESS:             "<",
	LESS_OR_EQUAL:    "<=",
	GREATER:          ">",
	GREATER_OR_EQUAL: ">=",
	LOGICAL_AND:      "&&",
	LOGICAL_OR:       "||",

	SHOW:           "show",
	DATABASES:      "databases",
	TABLES:         "tables",
	VARIABLES:      "variables",
	IF:             "if",
	NOT:            "not",
	EXISTS:         "exists",
	CREATE:         "create",
	DATABASE:       "database",
	TABLE:          "table",
	INDEX:          "index",
	UNIQUE:         "unique",
	SPATIAL:        "spatial",
	FULLTEXT:       "fulltext",
	DROP:           "drop",
	INSERT:         "insert",
	DELETE:         "delete",
	UPDATE:         "update",
	SELECT:         "select",
	FROM:           "from",
	WHERE:          "where",
	GROUP:          "group",
	HAVING:         "having",
	ORDER:          "order",
	LIMIT:          "limit",
	AS:             "as",
	SET:            "set",
	BY:             "by",
	AES:            "aes",
	DESC:           "desc",
	INTO:           "into",
	VALUES:         "values",
	INNER:          "inner",
	LEFT:           "left",
	RIGHT:          "right",
	JOIN:           "join",
	ON:             "on",
	IN:             "in",
	AND:            "and",
	OR:             "or",
	BETWEEN:        "between",
	LIKE:           "like",
	DISTINCT:       "distinct",
	UNION:          "union",
	ALL:            "all",
	TRUNCATE:       "truncate",
	PRIMARY:        "primary",
	KEY:            "key",
	UNSIGNED:       "unsigned",
	ZEROFILL:       "zerofill",
	AUTO_INCREMENT: "auto_increment",
	DEFAULT:        "default",
	COLUMN_COMMENT: "comment",
	TINYINT:        "tinyint",
	SMALLINT:       "smallint",
	MEDIUMINT:      "mediumint",
	INT:            "int",
	INTEGER:        "integer",
	BIGINT:         "bigint",
	FLOAT:          "float",
	DOUBLE:         "double",
	DECIMAL:        "decimal",
	DATE:           "date",
	TIME:           "time",
	DATETIME:       "datetime",
	TIMESTAMP:      "timestamp",
	CHAR:           "char",
	VARCHAR:        "varchar",
	BINARY:         "binary",
	VARBINARY:      "varbinary",
	TEXT:           "text",
	BLOB:           "blob",
}

var keywordMap = map[string]Token{
	"true":           BOOLEAN,
	"false":          BOOLEAN,
	"null":           NULL,
	"show":           SHOW,
	"databases":      DATABASES,
	"tables":         TABLES,
	"variables":      VARIABLES,
	"if":             IF,
	"not":            NOT,
	"exists":         EXISTS,
	"create":         CREATE,
	"database":       DATABASE,
	"table":          TABLE,
	"index":          INDEX,
	"unique":         UNIQUE,
	"spatial":        SPATIAL,
	"fulltext":       FULLTEXT,
	"drop":           DROP,
	"insert":         INSERT,
	"delete":         DELETE,
	"update":         UPDATE,
	"select":         SELECT,
	"from":           FROM,
	"where":          WHERE,
	"group":          GROUP,
	"having":         HAVING,
	"order":          ORDER,
	"limit":          LIMIT,
	"as":             AS,
	"set":            SET,
	"by":             BY,
	"aes":            AES,
	"desc":           DESC,
	"into":           INTO,
	"values":         VALUES,
	"inner":          INNER,
	"left":           LEFT,
	"right":          RIGHT,
	"join":           JOIN,
	"on":             ON,
	"in":             IN,
	"and":            AND,
	"or":             OR,
	"between":        BETWEEN,
	"like":           LIKE,
	"distinct":       DISTINCT,
	"union":          UNION,
	"all":            ALL,
	"truncate":       TRUNCATE,
	"primary":        PRIMARY,
	"key":            KEY,
	"unsigned":       UNSIGNED,
	"zerofill":       ZEROFILL,
	"auto_increment": AUTO_INCREMENT,
	"default":        DEFAULT,
	"comment":        COLUMN_COMMENT,
	"tinyint":        TINYINT,
	"smallint":       SMALLINT,
	"mediumint":      MEDIUMINT,
	"int":            INT,
	"integer":        INTEGER,
	"bigint":         BIGINT,
	"float":          FLOAT,
	"double":         DOUBLE,
	"decimal":        DECIMAL,
	"date":           DATE,
	"time":           TIME,
	"datetime":       DATETIME,
	"timestamp":      TIMESTAMP,
	"char":           CHAR,
	"varchar":        VARCHAR,
	"binary":         BINARY,
	"varbinary":      VARBINARY,
	"text":           TEXT,
	"blob":           BLOB,
}

func IsKeyword(k string) (Token, bool) {
	tkn, exists := keywordMap[k]
	return tkn, exists
}

var fieldTypeMap = map[Token]byte{
	TINYINT:   common.FIELD_TYPE_TINY,
	SMALLINT:  common.FIELD_TYPE_SHORT,
	MEDIUMINT: common.FIELD_TYPE_INT24,
	INT:       common.FIELD_TYPE_LONG,
	INTEGER:   common.FIELD_TYPE_LONG,
	BIGINT:    common.FIELD_TYPE_LONGLONG,
	FLOAT:     common.FIELD_TYPE_FLOAT,
	DOUBLE:    common.FIELD_TYPE_DOUBLE,
	DECIMAL:   common.FIELD_TYPE_DECIMAL,
	DATE:      common.FIELD_TYPE_DATE,
	TIME:      common.FIELD_TYPE_TIME,
	DATETIME:  common.FIELD_TYPE_DATETIME,
	TIMESTAMP: common.FIELD_TYPE_TIMESTAMP,
	CHAR:      common.FIELD_TYPE_STRING,
	VARCHAR:   common.FIELD_TYPE_VARCHAR,
	BINARY:    common.FIELD_TYPE_VARCHAR,
	VARBINARY: common.FIELD_TYPE_VARCHAR,
	TEXT:      common.FIELD_TYPE_BLOB,
	BLOB:      common.FIELD_TYPE_BLOB,
}

func GetFieldType(tkn Token) byte {
	return fieldTypeMap[tkn]
}
