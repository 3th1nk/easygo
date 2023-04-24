package mathUtil

import "testing"

func TestMaxInt64(t *testing.T) {
	if n := MaxInt64(1, 2, 3, 4, 5); n != 5 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MaxInt64(1, 2); n != 2 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MaxInt64(1); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
}

func TestMaxInt32(t *testing.T) {
	if n := MaxInt32(1, 2, 3, 4, 5); n != 5 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MaxInt32(1, 2); n != 2 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MaxInt32(1); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
}

func TestMaxInt(t *testing.T) {
	if n := MaxInt(1, 2, 3, 4, 5); n != 5 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MaxInt(1, 2); n != 2 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MaxInt(1); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
}

func TestMinInt64(t *testing.T) {
	if n := MinInt64(1, 2, 3, 4, 5); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MinInt64(1, 2); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MinInt64(1); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
}

func TestMinInt32(t *testing.T) {
	if n := MinInt32(1, 2, 3, 4, 5); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MinInt32(1, 2); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MinInt32(1); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
}

func TestMinInt(t *testing.T) {
	if n := MinInt(1, 2, 3, 4, 5); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MinInt(1, 2); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
	if n := MinInt(1); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
}
