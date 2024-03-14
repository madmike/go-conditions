package conditions

type SimpleOperatorsEnum string

const (
	NULL      SimpleOperatorsEnum = "$null"      // Represents the null operator
	DEFINED   SimpleOperatorsEnum = "$defined"   // Represents the defined operator
	UNDEFINED SimpleOperatorsEnum = "$undefined" // Represents the undefined operator
	EXIST     SimpleOperatorsEnum = "$exist"     // Represents the exist operator
	EMPTY     SimpleOperatorsEnum = "$empty"     // Represents the empty operator
	BLANK     SimpleOperatorsEnum = "$blank"     // Represents the blank operator
	TRULY     SimpleOperatorsEnum = "$truly"     // Represents the truly operator
	FALSY     SimpleOperatorsEnum = "$falsy"     // Represents the falsy operator
)

type CommonOperatorsEnum string

const (
	EQ      CommonOperatorsEnum = "$eq"      // Represents the equal operator
	NE      CommonOperatorsEnum = "$ne"      // Represents the not equal operator
	LT      CommonOperatorsEnum = "$lt"      // Represents the less than operator
	GT      CommonOperatorsEnum = "$gt"      // Represents the greater than operator
	LTE     CommonOperatorsEnum = "$lte"     // Represents the less than or equal to operator
	GTE     CommonOperatorsEnum = "$gte"     // Represents the greater than or equal to operator
	RE      CommonOperatorsEnum = "$re"      // Represents the regular expression operator
	IN      CommonOperatorsEnum = "$in"      // Represents the in operator
	NI      CommonOperatorsEnum = "$ni"      // Represents the not in operator
	SW      CommonOperatorsEnum = "$sw"      // Represents the starts with operator
	EW      CommonOperatorsEnum = "$ew"      // Represents the ends with operator
	INCL    CommonOperatorsEnum = "$incl"    // Represents the includes operator
	EXCL    CommonOperatorsEnum = "$excl"    // Represents the excludes operator
	HAS     CommonOperatorsEnum = "$has"     // Represents the has operator
	POWER   CommonOperatorsEnum = "$power"   // Represents the power operator
	BETWEEN CommonOperatorsEnum = "$between" // Represents the between operator
	SOME    CommonOperatorsEnum = "$some"    // Represents the some operator
	EVERY   CommonOperatorsEnum = "$every"   // Represents the every operator
	NOONE   CommonOperatorsEnum = "$noone"   // Represents the no one operator
)

type LogicOperatorsEnum string

const (
	OR  LogicOperatorsEnum = "$or"  // Represents the logical OR operator
	XOR LogicOperatorsEnum = "$xor" // Represents the logical XOR operator
	AND LogicOperatorsEnum = "$and" // Represents the logical AND operator
	NOT LogicOperatorsEnum = "$not" // Represents the logical NOT operator
)

var stringToSimpleOperator = map[string]SimpleOperatorsEnum{
	"$null":      NULL,
	"$defined":   DEFINED,
	"$undefined": UNDEFINED,
	"$empty":     EMPTY,
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
	"$in":      IN,
	"$ni":      NI,
	"$re":      RE,
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
