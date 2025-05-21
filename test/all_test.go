package test

import (
	"testing"

	"github.com/vedadiyan/sqlparser/pkg/sqlparser"
)

func TestAll(t *testing.T) {
	query := "SELECT * FROM X PARALLEL INNER JOIN Y ON X.id = Y.id"

	parser := sqlparser.Parser{}
	parsed, err := parser.Parse(query)
	if err != nil {
		t.FailNow()
	}
	jt := parsed.(*sqlparser.Select).From[0].(*sqlparser.JoinTableExpr).Join
	_ = jt
	_ = parsed
}
