package squirrel

import "testing"

func TestQuestion(t *testing.T) {
	s := Question.String()
	expect := "?"
	if s != expect {
		t.Errorf("expected %v, got %v", expect, s)
	}

	sql := "x = ? AND y = ?"
	s = Question.ReplacePlaceholders(sql)
	if s != sql {
		t.Errorf("expected %v, got %v", sql, s)
	}
}

func TestDollar(t *testing.T) {
	s := Dollar.String()
	expect := "$?"
	if s != expect {
		t.Errorf("expected %v, got %v", expect, s)
	}

	sql := "x = $? AND y = $?"
	s = Dollar.ReplacePlaceholders(sql)
	expect = "x = $1 AND y = $2"
	if s != expect {
		t.Errorf("expected %v, got %v", expect, s)
	}
}
