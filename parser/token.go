package parser

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

	IF       // if
	NOT      // not
	EXISTS   // exists
	CREATE   // create
	DATABASE // database
	TABLE    // table
	INDEX    // index
	UNIQUE   // unique
	SPATIAL  // spatial
	FULLTEXT // fulltext
	DROP     // drop
	INTO     // into
	VALUES   // values
	INSERT   // insert
	DELETE   // delete
	UPDATE   // update
	SELECT   // select
	FROM     // from
	WHERE    // where
	GROUP    // group
	HAVING   // having
	ORDER    // order
	LIMIT    // limit
	AS       // as
	BY       // by
	AES      // aes
	DESC     // desc
	INNER    // inner
	LEFT     // left
	RIGHT    // right
	JOIN     // join
	ON       // on
	IN       // in
	AND      // and
	OR       // or
	BETWEEN  // between
	LIKE     // like
	DISTINCT // distinct
	UNION    // union
	ALL      // all
	TRUNCATE // truncate
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

	IF:       "if",
	NOT:      "not",
	EXISTS:   "exists",
	CREATE:   "create",
	DATABASE: "database",
	TABLE:    "table",
	INDEX:    "index",
	UNIQUE:   "unique",
	SPATIAL:  "spatial",
	FULLTEXT: "fulltext",
	DROP:     "drop",
	INSERT:   "insert",
	DELETE:   "delete",
	UPDATE:   "update",
	SELECT:   "select",
	FROM:     "from",
	WHERE:    "where",
	GROUP:    "group",
	HAVING:   "having",
	ORDER:    "order",
	LIMIT:    "limit",
	AS:       "as",
	BY:       "by",
	AES:      "aes",
	DESC:     "desc",
	INTO:     "into",
	VALUES:   "values",
	INNER:    "inner",
	LEFT:     "left",
	RIGHT:    "right",
	JOIN:     "join",
	ON:       "on",
	IN:       "in",
	AND:      "and",
	OR:       "or",
	BETWEEN:  "between",
	LIKE:     "like",
	DISTINCT: "distinct",
	UNION:    "union",
	ALL:      "all",
	TRUNCATE: "truncate",
}

var keywordMap = map[string]Token{
	"true":     BOOLEAN,
	"false":    BOOLEAN,
	"null":     NULL,
	"if":       IF,
	"not":      NOT,
	"exists":   EXISTS,
	"create":   CREATE,
	"database": DATABASE,
	"table":    TABLE,
	"index":    INDEX,
	"unique":   UNIQUE,
	"spatial":  SPATIAL,
	"fulltext": FULLTEXT,
	"drop":     DROP,
	"insert":   INSERT,
	"delete":   DELETE,
	"update":   UPDATE,
	"select":   SELECT,
	"from":     FROM,
	"where":    WHERE,
	"group":    GROUP,
	"having":   HAVING,
	"order":    ORDER,
	"limit":    LIMIT,
	"as":       AS,
	"by":       BY,
	"aes":      AES,
	"desc":     DESC,
	"into":     INTO,
	"values":   VALUES,
	"inner":    INNER,
	"left":     LEFT,
	"right":    RIGHT,
	"join":     JOIN,
	"on":       ON,
	"in":       IN,
	"and":      AND,
	"or":       OR,
	"between":  BETWEEN,
	"like":     LIKE,
	"distinct": DISTINCT,
	"union":    UNION,
	"all":      ALL,
	"truncate": TRUNCATE,
}

func IsKeyword(k string) (Token, bool) {
	tkn, exists := keywordMap[k]
	return tkn, exists
}
