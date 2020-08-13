// +build unittest

package mycrons

import "testing"

func TestParser_Parse(t *testing.T) {
	var err error
	var s Scheduler

	//----------------- week 正常情况 --------------------- //
	//
	if s, err = NewParser().Parse("week 30 10 1"); err != nil {
		t.Fatalf("failed, %v", err)
	}
	if _, ok := s.(*weekScheduler); !ok {
		t.Fatal("bad type of scheduler returned")
	}
	//
	if s, err = NewParser().Parse("week 0,30,45 10,20 1,3,4,7"); err != nil {
		t.Fatalf("failed, %v", err)
	}
	if _, ok := s.(*weekScheduler); !ok {
		t.Fatal("bad type of scheduler returned")
	}
	//
	if s, err = NewParser().Parse("week 0,30,45 10,20 1-7"); err != nil {
		t.Fatalf("failed, %v", err)
	}
	if _, ok := s.(*weekScheduler); !ok {
		t.Fatal("bad type of scheduler returned")
	}

	// ----------------- month 正常情况 ------------------- //
	//
	if s, err = NewParser().Parse("month 30 10 1"); err != nil {
		t.Fatalf("failed, %v", err)
	}
	if _, ok := s.(*monthScheduler); !ok {
		t.Fatal("bad type of scheduler returned")
	}
	//
	if s, err = NewParser().Parse("month 0,30,45 10,20 1,31"); err != nil {
		t.Fatalf("failed, %v", err)
	}
	if _, ok := s.(*monthScheduler); !ok {
		t.Fatal("bad type of scheduler returned")
	}
	//
	if s, err = NewParser().Parse("month 0,30,45 10,20 1-31"); err != nil {
		t.Fatalf("failed, %v", err)
	}
	if _, ok := s.(*monthScheduler); !ok {
		t.Fatal("bad type of scheduler returned")
	}

	// ------------------- week 异常情况 --------------------- //
	if _, err = NewParser().Parse("wee 30 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse(" 30 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 *"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 * 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week * 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week -30 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 -10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 -1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week A30 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 A10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 A1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week   30       10           1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 60 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 24 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 0"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 8"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 1-7-8"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 1--8"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 1-8"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("week 30 10 1-"); err == nil {
		t.Fatal("failed")
	}

	// ------------------- month 异常情况 --------------------- //
	if _, err = NewParser().Parse("mon 30 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse(" 30 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 *"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 * 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month * 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month -30 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 -10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 -1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month A30 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 A10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 A1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month   30       10           1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 60 10 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 24 1"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 0"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 32"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 1-7-8"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 1--8"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 1-32"); err == nil {
		t.Fatal("failed")
	}
	if _, err = NewParser().Parse("month 30 10 1-"); err == nil {
		t.Fatal("failed")
	}
}
