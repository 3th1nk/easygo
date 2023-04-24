package timeUtil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestJsonTime_Json(t *testing.T) {
	now := ToJsonTime(time.Now())
	objA := &JsonTimeTest{
		Date:       now,
		Time:       now,
		Timestamp:  now,
		Timestamp2: now,
	}
	objB := &JsonTimeTest{}

	dataA, _ := json.Marshal(objA)
	json.Unmarshal(dataA, objB)
	dataB, _ := json.Marshal(objB)
	strA, strB := string(dataA), string(dataB)

	if strA != strB {
		t.Error(fmt.Sprintf("assert faild: expect %v, but %v", strA, strB))
	}

	format := "2006-01-02 15:04:05.999999"
	if !objA.Date.Equal(objB.Date.Time) {
		t.Error(fmt.Sprintf("assert faild: expect %v, but %v", objA.Date.Format(format), objB.Date.Format(format)))
	}
	if !objA.Time.Equal(objB.Time.Time) {
		t.Error(fmt.Sprintf("assert faild: expect %v, but %v", objA.Time.Format(format), objB.Time.Format(format)))
	}
	if !objA.Timestamp.Equal(objB.Timestamp.Time) {
		t.Error(fmt.Sprintf("assert faild: expect %v, but %v", objA.Timestamp.Format(format), objB.Timestamp.Format(format)))
	}
	if !objA.Timestamp2.Equal(objB.Timestamp2.Time) {
		t.Error(fmt.Sprintf("assert faild: expect %v, but %v", objA.Timestamp2.Format(format), objB.Timestamp2.Format(format)))
	}
}

// ------------------------------------------------------------------------------ test.json_time_test
type JsonTimeTest struct {
	Id         int      `json:"id"`
	Date       JsonTime `json:"date,omitempty"`
	Time       JsonTime `json:"time,omitempty"`
	Timestamp  JsonTime `json:"timestamp,omitempty"`
	Timestamp2 JsonTime `json:"timestamp2,omitempty"`
}

func TestJsonTime_Equal(t *testing.T) {
	a := &JsonTimeTest{Id: 1}
	b := &JsonTimeTest{Id: 1}
	t.Log(a == b)
	t.Log(reflect.DeepEqual(a, b))

	c := &JsonTimeTest{Id: 2, Date: ToJsonTime(BeginOfDay(time.Now()))}
	d := &JsonTimeTest{Id: 2, Date: ToJsonTime(BeginOfDay(time.Now()))}
	t.Log(c == d)
	t.Log(reflect.DeepEqual(c, d))
}
