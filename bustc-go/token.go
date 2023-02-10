package main

type TokenKind int

const (
	InvalidToken = iota
	IdentifierToken
	NumberToken
	RuneToken
	StringToken
	OpenBraceToken
	CloseBraceToken
	OpenParenthesisToken
	CloseParenthesisToken
	OpenBracketToken
	CloseBracketToken
	FuncToken
	PackageToken
	TypeToken
	StructToken
	EnumToken
	ReturnToken
	IfToken
	ElseToken
	ForToken
	BreakToken
	ContinueToken
	InitToken
	AssignToken
	ScopeToken
	AddToken
	SubToken
	MulToken
	DivToken
	AndToken
	EqualToken
	LessToken
	GreaterToken
	BoolAndToken
	BoolOrToken
	CommaToken
	LineCommentToken
)

func (k TokenKind) String() string {
	switch k {
	case InvalidToken:
		return "[invalid]"
	case IdentifierToken:
		return "Identifier"
	case NumberToken:
		return "Number"
	case RuneToken:
		return "Rune"
	case StringToken:
		return "String"
	case OpenBraceToken:
		return "OpenBrace"
	case CloseBraceToken:
		return "CloseBrace"
	case OpenParenthesisToken:
		return "OpenParenthesis"
	case CloseParenthesisToken:
		return "CloseParenthesis"
	case OpenBracketToken:
		return "OpenBracket"
	case CloseBracketToken:
		return "CloseBracket"
	case TypeToken:
		return "TypeToken"
	case FuncToken:
		return "FunctionToken"
	case PackageToken:
		return "PackageToken"
	case StructToken:
		return "StructToken"
	case EnumToken:
		return "EnumToken"
	case ReturnToken:
		return "ReturnToken"
	case IfToken:
		return "IfToken"
	case ElseToken:
		return "ElseToken"
	case ForToken:
		return "ForToken"
	case BreakToken:
		return "BreakToken"
	case ContinueToken:
		return "ContinueToken"
	case InitToken:
		return "InitToken"
	case AssignToken:
		return "AssignToken"
	case ScopeToken:
		return "ScopeToken"
	case AddToken:
		return "AddToken"
	case SubToken:
		return "SubToken"
	case MulToken:
		return "MulToken"
	case DivToken:
		return "DivToken"
	case AndToken:
		return "AndToken"
	case EqualToken:
		return "EqualToken"
	case LessToken:
		return "LessToken"
	case GreaterToken:
		return "GreaterToken"
	case BoolAndToken:
		return "BoolAndToken"
	case BoolOrToken:
		return "BoolOrToken"
	case CommaToken:
		return "CommaToken"
	case LineCommentToken:
		return "LineCommentToken"
	default:
		return "[unknown]"
	}
}

type Token struct {
	Kind     TokenKind
	Text     string
	Position Position
}

var keywords = map[string]TokenKind{
	"type":     TypeToken,
	"func":     FuncToken,
	"package":  PackageToken,
	"struct":   StructToken,
	"enum":     EnumToken,
	"return":   ReturnToken,
	"if":       IfToken,
	"else":     ElseToken,
	"for":      ForToken,
	"break":    BreakToken,
	"continue": ContinueToken,
}
