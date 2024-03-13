package conditions

type SimpleOperatorsEnum string

const (
	NULL      SimpleOperatorsEnum = "$null"
	DEFINED   SimpleOperatorsEnum = "$defined"
	UNDEFINED SimpleOperatorsEnum = "$undefined"
	EXIST     SimpleOperatorsEnum = "$exist"
	EMPTY     SimpleOperatorsEnum = "$empty"
	BLANK     SimpleOperatorsEnum = "$blank"
	TRULY     SimpleOperatorsEnum = "$truly"
	FALSY     SimpleOperatorsEnum = "$falsy"
)

type CommonOperatorsEnum string

const (
	EQ      CommonOperatorsEnum = "$eq"
	NE      CommonOperatorsEnum = "$ne"
	LT      CommonOperatorsEnum = "$lt"
	GT      CommonOperatorsEnum = "$gt"
	LTE     CommonOperatorsEnum = "$lte"
	GTE     CommonOperatorsEnum = "$gte"
	RE      CommonOperatorsEnum = "$re"
	IN      CommonOperatorsEnum = "$in"
	NI      CommonOperatorsEnum = "$ni"
	SW      CommonOperatorsEnum = "$sw"
	EW      CommonOperatorsEnum = "$ew"
	INCL    CommonOperatorsEnum = "$incl"
	EXCL    CommonOperatorsEnum = "$excl"
	HAS     CommonOperatorsEnum = "$has"
	POWER   CommonOperatorsEnum = "$power"
	BETWEEN CommonOperatorsEnum = "$between"
	SOME    CommonOperatorsEnum = "$some"
	EVERY   CommonOperatorsEnum = "$every"
	NOONE   CommonOperatorsEnum = "$noone"
)

type LogicOperatorsEnum string

const (
	OR  LogicOperatorsEnum = "$or"
	XOR LogicOperatorsEnum = "$xor"
	AND LogicOperatorsEnum = "$and"
	NOT LogicOperatorsEnum = "$not"
)

var stringToSimpleOperator = map[string]SimpleOperatorsEnum{
	"$null":      NULL,
	"$defined":   DEFINED,
	"$undefined": UNDEFINED,
	"$exist":     EXIST,
	"$blank":     BLANK,
	"$truly":     TRULY,
	"$falsy":     FALSY,
}

var stringToCommonOperator = map[string]CommonOperatorsEnum{
	"$eq":      EQ,
	"$ne":      NE,
	"$lt":      LT,
	"$gt":      GT,
	"$lte":     LTE,
	"$gte":     GTE,
	"$re":      RE,
	"$in":      IN,
	"$ni":      NI,
	"$sw":      SW,
	"$ew":      EW,
	"$incl":    INCL,
	"$excl":    EXCL,
	"$has":     HAS,
	"$power":   POWER,
	"$between": BETWEEN,
	"$some":    SOME,
	"$every":   EVERY,
	"$noone":   NOONE,
}

var stringToLogicOperator = map[string]LogicOperatorsEnum{
	"$or":  OR,
	"$xor": XOR,
	"$and": AND,
	"$not": NOT,
}

var SimpleOperators = []SimpleOperatorsEnum{NULL, DEFINED, UNDEFINED, EXIST, EMPTY, BLANK, TRULY, FALSY}
var CommonOperators = []CommonOperatorsEnum{EQ, NE, LT, GT, LTE, GTE, RE, IN, NI, SW, EW, INCL, EXCL, HAS, POWER, BETWEEN, SOME, EVERY, NOONE}
var LogicOperators = []LogicOperatorsEnum{OR, XOR, AND, NOT}
