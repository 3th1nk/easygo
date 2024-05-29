package influxdb

import (
	"testing"
)

func TestExpr(t *testing.T) {
	t.Log(RawExpr("n = 1"))
	t.Log(RawExpr("s = 'a'"))
	t.Log(RawExpr("time >= '2024-01-01'"))
	t.Log(RawExpr("time >= '2024-01-01' AND time <= '2024-05-01'"))
}

func TestSimpleExpr(t *testing.T) {
	t.Log(Expr("n", "=", 1))
	t.Log(Expr("s", "=", "a"))
	t.Log(Expr("s", "=", "a AND b != \"c\""))
	t.Log(Expr("time", ">=", "2024-01-01"))
	t.Log(Expr("time", ">=", 1716355881))
}

func TestBetween(t *testing.T) {
	t.Log(Between("n", 1, 2))
	t.Log(Between("n", 1, nil))
	t.Log(Between("n", nil, 2))

	t.Log(Between("time", "2024-01-01", "2024-05-01"))
	t.Log(Between("time", "2024-01-01", nil))
	t.Log(Between("time", nil, "2024-05-01"))
}

func TestAnd(t *testing.T) {
	t.Log(And(RawExpr("n = 1"), RawExpr("s = 'a'")))
	t.Log(And(RawExpr("n = 1"), And(RawExpr("s = 'a'"), RawExpr("time >= '2024-01-01'"))))
	t.Log(And(Expr("n", "=", 1), Expr("s", "=", "a")))
}

func TestOr(t *testing.T) {
	t.Log(Or(RawExpr("n = 1"), RawExpr("s = 'a'")))
	t.Log(Or(RawExpr("n = 1"), Or(RawExpr("s = 'a'"), RawExpr("time >= '2024-01-01'"))))
}

func TestAndOr(t *testing.T) {
	t.Log(And(RawExpr("n = 1"), Or(RawExpr("s = 'a'"), RawExpr("time >= '2024-01-01'"))))
	t.Log(Or(RawExpr("n = 1"), And(RawExpr("s = 'a'"), RawExpr("time >= '2024-01-01'"))))
}

func TestIn(t *testing.T) {
	t.Log(In("n", 1, 2, 3))
	t.Log(In("s", "a", "b", "c"))
	t.Log(In("time", "2024-01-01", "2024-05-01"))
}

func TestNotIn(t *testing.T) {
	t.Log(NotIn("n", 1, 2, 3))
	t.Log(NotIn("s", "a", "b", "c"))
	t.Log(NotIn("time", "2024-01-01", "2024-05-01"))
}

func TestLike(t *testing.T) {
	t.Log(Match("s", "a.*"))
	t.Log(Match("s", "^a.*"))
	t.Log(Match("s", "a.*$"))
}

func TestNotLike(t *testing.T) {
	t.Log(NotMatch("s", "a.*"))
	t.Log(NotMatch("s", "^a.*"))
	t.Log(NotMatch("s", "a.*$"))
}
