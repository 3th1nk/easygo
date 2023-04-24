package arrUtil

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func testIntData(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(n)
	}
	return arr
}

/*
=== RUN   TestPerform
    distinct_perform_test.go:56: ([81 87 47 59 81 18 25 40 56 0 94 11 62 89 28 74 11 45 37 6 95 66 28 58 47 47 87 88 90 15 41 8 87 31 29 56 37 31 85 26 13 90 94 63 33 47 78 24 59 53 57 21 89 99 0 5 88 38 3 55 51 10 5 56 66 28 61 2 83 46 63 76 2 18 47 94 77 63 96 20 23 53 37 33 41 59 33 43 91 2 78 36 46 7 40 3 52 43 5 98])
        distinctIntWithArr([81 87 47 59 18 25 40 56 0 94 11 62 89 28 74 45 37 6 95 66 58 88 90 15 41 8 31 29 85 26 13 63 33 78 24 53 57 21 99 5 38 3 55 51 10 61 2 83 46 76 77 96 20 23 43 91 36 7 52 98])
        distinctIntWithMap([81 87 47 59 18 25 40 56 0 94 11 62 89 28 74 45 37 6 95 66 58 88 90 15 41 8 31 29 85 26 13 63 33 78 24 53 57 21 99 5 38 3 55 51 10 61 2 83 46 76 77 96 20 23 43 91 36 7 52 98])
    distinct_perform_test.go:86: perform(distinctIntWithMap) count: 200, took: 38.998µs
    distinct_perform_test.go:86: perform(distinctIntWithArr) count: 200, took: 53.611µs
    distinct_perform_test.go:86: perform(distinctIntWithMap) count: 300, took: 60.211µs
    distinct_perform_test.go:86: perform(distinctIntWithArr) count: 300, took: 102.851µs
    distinct_perform_test.go:86: perform(distinctIntWithMap) count: 400, took: 109.276µs
    distinct_perform_test.go:86: perform(distinctIntWithArr) count: 400, took: 146.587µs
    distinct_perform_test.go:86: perform(distinctIntWithMap) count: 500, took: 93.193µs
    distinct_perform_test.go:86: perform(distinctIntWithArr) count: 500, took: 251.563µs
    distinct_perform_test.go:86: perform(distinctIntWithMap) count: 1000, took: 166.41µs
    distinct_perform_test.go:86: perform(distinctIntWithArr) count: 1000, took: 752.213µs
    distinct_perform_test.go:86: perform(distinctIntWithMap) count: 2000, took: 296.376µs
    distinct_perform_test.go:86: perform(distinctIntWithArr) count: 2000, took: 2.80948ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithMap) count: 200, took: 23.092413ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithArr) count: 200, took: 13.118421ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithMap) count: 300, took: 36.16878ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithArr) count: 300, took: 34.677486ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithMap) count: 400, took: 34.03811ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithArr) count: 400, took: 54.002012ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithMap) count: 500, took: 36.782406ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithArr) count: 500, took: 86.124142ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithMap) count: 1000, took: 103.779885ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithArr) count: 1000, took: 317.911776ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithMap) count: 2000, took: 113.741494ms
    distinct_perform_test.go:101: go(1000) perform(distinctIntWithArr) count: 2000, took: 1.270327897s
--- PASS: TestPerform (2.13s)
PASS
*/
func TestPerform(t *testing.T) {
	arr := testIntData(100)
	da := distinctIntWithArr(arr)
	dm := distinctIntWithMap(arr)

	for i, v := range da {
		if dm[i] != v || len(da) != len(dm) {
			t.Fatalf("(%v) \ndistinctIntWithArr(%d.%v) \ndistinctIntWithMap(%d.%v)", arr, len(da), da, len(dm), dm)
		}
	}
	t.Logf("(%v) \ndistinctIntWithArr(%v)  \ndistinctIntWithMap(%v)", arr, da, dm)

	for _, n := range []int{
		200, 300, 400, 500, 1000,
	} {
		arr := testIntData(n)
		testPerformWithSinge(t, "distinctIntWithMap", n, func() {
			distinctIntWithMap(arr)
		})
		testPerformWithSinge(t, "distinctIntWithArr", n, func() {
			distinctIntWithArr(arr)
		})
	}

	for _, n := range []int{
		200, 300, 400, 500, 1000,
	} {
		arr := testIntData(n)
		testPerformWithGo(t, "distinctIntWithMap", n, func() {
			distinctIntWithMap(arr)
		})
		testPerformWithGo(t, "distinctIntWithArr", n, func() {
			distinctIntWithArr(arr)
		})
	}
}

func testPerformWithSinge(t *testing.T, fName string, n int, f func()) {
	start := time.Now()
	f()
	t.Logf("perform(%s) count: %d, took: %s", fName, n, time.Since(start).String())
}

func testPerformWithGo(t *testing.T, fName string, n int, f func()) {
	start := time.Now()
	wg := sync.WaitGroup{}
	goNum := 1000
	for i := 0; i < goNum; i++ {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}
	wg.Wait()
	t.Logf("go(%d) perform(%s) count: %d, took: %s", goNum, fName, n, time.Since(start).String())
}

// 产生重复的数据
func testStringData(n int) []string {

	arrN := n / 3
	strArr := make([]string, arrN)
	for i := 0; i < arrN; i++ {
		strArr[i] = strUtil.RandULN(10)
	}

	arr := make([]string, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			arr[i] = strings.ToLower(strArr[rand.Intn(arrN)])
		}
		arr[i] = strArr[rand.Intn(arrN)]
	}

	return arr
}

/*
=== RUN   TestPerformDistinctString
    distinct_perform_test.go:186: ([EChY1vOYTm 31EID5ZuEk MoMmTFAD7Q aRboFdqMFH v19QklBFil Bj9PhlNHIW 0oPC82XxoG lkyA9Txdtj cVxd5Jb0Ih FHrVoKIh8d LsCIbDNMyI hAtAMZyUEv hAtAMZyUEv 55oL6ksgA7 eLWTBDzkDX xOiHWWDGHs lkyA9Txdtj bu66k9wV9u RR8jFcsltH KwcVdVhj0i LsCIbDNMyI hAtAMZyUEv JnwR65pbdU cVxd5Jb0Ih Bj9PhlNHIW EChY1vOYTm Bj9PhlNHIW Bj9PhlNHIW aRboFdqMFH MoMmTFAD7Q bu66k9wV9u lkyA9Txdtj aRboFdqMFH aRboFdqMFH JnwR65pbdU OCClUyL8iK DilZznzNkl jyGZHxtYWA lkyA9Txdtj 5evlmtPk8E aRboFdqMFH OCClUyL8iK Bj9PhlNHIW cVxd5Jb0Ih 31EID5ZuEk cVxd5Jb0Ih bu66k9wV9u RR8jFcsltH 0oPC82XxoG EChY1vOYTm FHrVoKIh8d 9uFaf3ohr4 hAtAMZyUEv OCClUyL8iK lkyA9Txdtj 4M7MqvpdDk T4dQrs2DVY 0oPC82XxoG n1pfML2YOs 9uFaf3ohr4 lkyA9Txdtj KwcVdVhj0i YjcKFcZUdd xOiHWWDGHs YjcKFcZUdd HVlfEitGuk MUludL2t84 T4dQrs2DVY 0oPC82XxoG rmBzCjN9ft lkyA9Txdtj 4M7MqvpdDk 9uFaf3ohr4 KwcVdVhj0i MoMmTFAD7Q n1pfML2YOs FHrVoKIh8d hAtAMZyUEv MUludL2t84 KwcVdVhj0i MoMmTFAD7Q LsCIbDNMyI bu66k9wV9u JnwR65pbdU HVlfEitGuk FHrVoKIh8d LsCIbDNMyI RR8jFcsltH 55oL6ksgA7 8xM6eGZnZx v19QklBFil 0oPC82XxoG KwcVdVhj0i RR8jFcsltH T4dQrs2DVY FHrVoKIh8d rmBzCjN9ft 0oPC82XxoG v19QklBFil rmBzCjN9ft])
        distinctStringWithArr([EChY1vOYTm 31EID5ZuEk MoMmTFAD7Q aRboFdqMFH v19QklBFil Bj9PhlNHIW 0oPC82XxoG lkyA9Txdtj cVxd5Jb0Ih FHrVoKIh8d LsCIbDNMyI hAtAMZyUEv 55oL6ksgA7 eLWTBDzkDX xOiHWWDGHs bu66k9wV9u RR8jFcsltH KwcVdVhj0i JnwR65pbdU OCClUyL8iK DilZznzNkl jyGZHxtYWA 5evlmtPk8E 9uFaf3ohr4 4M7MqvpdDk T4dQrs2DVY n1pfML2YOs YjcKFcZUdd HVlfEitGuk MUludL2t84 rmBzCjN9ft 8xM6eGZnZx])
        distinctStringWithMap([EChY1vOYTm 31EID5ZuEk MoMmTFAD7Q aRboFdqMFH v19QklBFil Bj9PhlNHIW 0oPC82XxoG lkyA9Txdtj cVxd5Jb0Ih FHrVoKIh8d LsCIbDNMyI hAtAMZyUEv 55oL6ksgA7 eLWTBDzkDX xOiHWWDGHs bu66k9wV9u RR8jFcsltH KwcVdVhj0i JnwR65pbdU OCClUyL8iK DilZznzNkl jyGZHxtYWA 5evlmtPk8E 9uFaf3ohr4 4M7MqvpdDk T4dQrs2DVY n1pfML2YOs YjcKFcZUdd HVlfEitGuk MUludL2t84 rmBzCjN9ft 8xM6eGZnZx])
    distinct_perform_test.go:86: perform(distinctStringWithArr) count: 100, took: 18.678µs
    distinct_perform_test.go:86: perform(distinctStringWithMap) count: 100, took: 86.269µs
    distinct_perform_test.go:86: perform(distinctStringWithArr) count: 150, took: 28.933µs
    distinct_perform_test.go:86: perform(distinctStringWithMap) count: 150, took: 32.149µs
    distinct_perform_test.go:86: perform(distinctStringWithArr) count: 180, took: 39.042µs
    distinct_perform_test.go:86: perform(distinctStringWithMap) count: 180, took: 23.62µs
    distinct_perform_test.go:86: perform(distinctStringWithArr) count: 200, took: 79.66µs
    distinct_perform_test.go:86: perform(distinctStringWithMap) count: 200, took: 24.204µs
    distinct_perform_test.go:86: perform(distinctStringWithArr) count: 300, took: 92.192µs
    distinct_perform_test.go:86: perform(distinctStringWithMap) count: 300, took: 57.223µs
    distinct_perform_test.go:86: perform(distinctStringWithArr) count: 500, took: 260.948µs
    distinct_perform_test.go:86: perform(distinctStringWithMap) count: 500, took: 162.177µs
    distinct_perform_test.go:86: perform(distinctStringWithArr) count: 1000, took: 1.090468ms
    distinct_perform_test.go:86: perform(distinctStringWithMap) count: 1000, took: 152.126µs
    distinct_perform_test.go:86: perform(distinctStringWithArr) count: 2000, took: 4.209385ms
    distinct_perform_test.go:86: perform(distinctStringWithMap) count: 2000, took: 284.612µs
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArr) count: 100, took: 5.703035ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMap) count: 100, took: 14.049943ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArr) count: 150, took: 19.267507ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMap) count: 150, took: 18.948897ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArr) count: 180, took: 18.796893ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMap) count: 180, took: 57.307523ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArr) count: 200, took: 19.576295ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMap) count: 200, took: 22.087175ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArr) count: 300, took: 44.766885ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMap) count: 300, took: 33.254446ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArr) count: 500, took: 101.339891ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMap) count: 500, took: 45.989781ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArr) count: 1000, took: 357.550102ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMap) count: 1000, took: 74.831492ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArr) count: 2000, took: 1.387933139s
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMap) count: 2000, took: 135.574696ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArrTrue) count: 100, took: 13.569381ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMapTrue) count: 100, took: 30.660014ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArrTrue) count: 150, took: 39.803813ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMapTrue) count: 150, took: 44.523137ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArrTrue) count: 180, took: 48.035068ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMapTrue) count: 180, took: 54.270795ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArrTrue) count: 200, took: 59.006433ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMapTrue) count: 200, took: 41.155885ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArrTrue) count: 300, took: 128.181085ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMapTrue) count: 300, took: 101.259211ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArrTrue) count: 500, took: 314.856462ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMapTrue) count: 500, took: 131.768043ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArrTrue) count: 1000, took: 1.334342458s
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMapTrue) count: 1000, took: 312.280217ms
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithArrTrue) count: 2000, took: 5.406671749s
    distinct_perform_test.go:101: go(1000) perform(distinctStringWithMapTrue) count: 2000, took: 468.982391ms
--- PASS: TestPerformDistinctString (10.90s)
PASS

*/

func TestPerformDistinctString(t *testing.T) {
	arr := testStringData(100)
	da := distinctStringWithArr(arr, true)
	dm := distinctStringWithMap(arr, true)
	for i, v := range da {
		if dm[i] != v || len(da) != len(dm) {
			t.Fatalf("(%s) \ndistinctStringWithArr(%s)  \ndistinctStringWithMap(%s)", arr, da, dm)
		}
	}
	t.Logf("(%s) \ndistinctStringWithArr(%s)  \ndistinctStringWithMap(%s)", arr, da, dm)

	for _, n := range []int{
		100, 200, 300, 400, 500, 1000,
	} {
		arr := testStringData(n)
		testPerformWithSinge(t, "distinctStringWithArr", n, func() {
			distinctStringWithArr(arr)
		})
		testPerformWithSinge(t, "distinctStringWithMap", n, func() {
			distinctStringWithMap(arr)
		})
	}

	for _, n := range []int{
		100, 200, 300, 400, 500, 1000,
	} {
		arr := testStringData(n)
		testPerformWithGo(t, "distinctStringWithArr", n, func() {
			distinctStringWithArr(arr)
		})
		testPerformWithGo(t, "distinctStringWithMap", n, func() {
			distinctStringWithMap(arr)
		})
	}

	for _, n := range []int{
		100, 200, 300, 400, 500, 1000,
	} {
		arr := testStringData(n)
		testPerformWithGo(t, "distinctStringWithArrTrue", n, func() {
			distinctStringWithArr(arr, true)
		})
		testPerformWithGo(t, "distinctStringWithMapTrue", n, func() {
			distinctStringWithMap(arr, true)
		})
	}
}
