package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nht1206/go-study/fq/list"
	"github.com/nht1206/go-study/fq/tokenizer"
)

type SignOp string

const (
	SignGt  SignOp = ">"
	SignLt  SignOp = "<"
	SignEq  SignOp = "="
	SignNeq SignOp = "!="
	SignGe  SignOp = ">="
	SignLe  SignOp = "<="
)

var SignOps = []SignOp{
	SignGt,
	SignLt,
	SignEq,
	SignNeq,
	SignGe,
	SignLe,
}

type JoinOp string

const (
	JoinAnd = JoinOp("&&")
	JoinOr  = JoinOp("||")
)

var JoinOps = []JoinOp{
	JoinAnd,
	JoinOr,
}

type Expr struct {
	Left  *tokenizer.Token
	Right *tokenizer.Token
	Op    SignOp
}

type ExprGroup struct {
	Item any
	Op   JoinOp
}

const (
	stepBeforeSign = iota
	stepSign
	stepAfterSign
	stepJoin
)

func Parse(filterQuery string) ([]ExprGroup, error) {
	s := tokenizer.NewTokenizer(strings.NewReader(filterQuery))

	result := make([]ExprGroup, 0)
	expr := Expr{}
	join := JoinAnd
	step := stepBeforeSign
	for s.HasNext() || s.Err != nil {
		t := s.Scan()

		if t.Type == tokenizer.TokenEOF {
			break
		}

		if t.Type == tokenizer.TokenWhiteSpace {
			continue
		}

		if t.Type == tokenizer.TokenUnsupported {
			return nil, errors.New("unexpected character found")
		}

		if t.Type == tokenizer.TokenGroup {
			groupResult, err := Parse(t.Literal)
			if err != nil {
				return nil, err
			}

			if len(groupResult) > 0 {
				result = append(result, ExprGroup{
					Item: groupResult,
					Op:   join,
				})
			}

			step = stepJoin
			continue
		}

		switch step {
		case stepBeforeSign:
			if t.Type != tokenizer.TokenText && t.Type != tokenizer.TokenNumber && t.Type != tokenizer.TokenIdentifier {
				return nil, fmt.Errorf("expected left operand (identifier, text or number), but got %q (%s)", t.Literal, t.Type)
			}
			expr.Left = t
			step = stepSign
			continue
		case stepSign:
			sign := SignOp(t.Literal)
			if t.Type != tokenizer.TokenOperator || !list.ExistInSlice(SignOps, sign) {
				return nil, fmt.Errorf("non deterministic sign operator %q", t.Literal)
			}
			expr.Op = sign
			step = stepAfterSign
			continue
		case stepAfterSign:
			if t.Type != tokenizer.TokenText && t.Type != tokenizer.TokenNumber && t.Type != tokenizer.TokenIdentifier {
				return nil, fmt.Errorf("expected right operand (identifier, text or number), but got %q (%s)", t.Literal, t.Type)
			}
			expr.Right = t

			result = append(result, ExprGroup{
				Op:   join,
				Item: expr,
			})

			step = stepJoin
			continue
		case stepJoin:
			join = JoinOp(t.Literal)
			if t.Type != tokenizer.TokenJoin || !list.ExistInSlice(JoinOps, join) {
				return nil, fmt.Errorf("non deterministic join operator %q (%s)", t.Literal, t.Type)
			}

			expr = Expr{}
			step = stepBeforeSign
			continue
		}
	}

	if step != stepJoin {
		return nil, errors.New("incomplete filter query")
	}

	return result, nil
}
