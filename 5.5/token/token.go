package token

type Token int

const (
	Invalid Token = iota

	Name
	Numeral
	LiteralString

	And
	Break
	Do
	Else
	Elseif
	End
	False
	For
	Function
	Global
	Goto
	If
	In
	Local
	Nil
	Not
	Or
	Repeat
	Return
	Then
	True
	Until
	While

	Plus
	Minus
	Asterisk
	Solidus
	Percent
	Caret
	Hash
	BitAnd
	Tilde
	BitOr
	DoubleLessThan
	DoubleGreaterThan
	DoubleSolidus
	DoubleEqual
	NotEqual
	LessThanEqual
	GreaterThanEqual
	LessThan
	GreaterThan
	Equal
	LeftParen
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket
	DoubleColon
	Semicolon
	Colon
	Comma
	Dot
	DoubleDot
	TripleDot
)

func (t Token) String() string {
	return names[t]
}

var names = []string{
	"Invalid",
	"Name",
	"Numeral",
	"LiteralString",
	"And",
	"Break",
	"Do",
	"Else",
	"Elseif",
	"End",
	"False",
	"For",
	"Function",
	"Global",
	"Goto",
	"If",
	"In",
	"Local",
	"Nil",
	"Not",
	"Or",
	"Repeat",
	"Return",
	"Then",
	"True",
	"Until",
	"While",
	"Plus",
	"Minus",
	"Asterisk",
	"Solidus",
	"Percent",
	"Caret",
	"Hash",
	"BitAnd",
	"Tilde",
	"BitOr",
	"DoubleLessThan",
	"DoubleGreaterThan",
	"DoubleSolidus",
	"DoubleEqual",
	"NotEqual",
	"LessThanEqual",
	"GreaterThanEqual",
	"LessThan",
	"GreaterThan",
	"Equal",
	"LeftParen",
	"RightParen",
	"LeftBrace",
	"RightBrace",
	"LeftBracket",
	"RightBracket",
	"DoubleColon",
	"Semicolon",
	"Colon",
	"Comma",
	"Dot",
	"DoubleDot",
	"TripleDot",
}
