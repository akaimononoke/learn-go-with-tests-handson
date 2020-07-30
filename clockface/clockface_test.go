package clockface

import (
	"math"
	"testing"
	"time"
)

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-7
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b Point) bool {
	return roughlyEqualFloat64(a.X, b.X) && roughlyEqualFloat64(a.Y, b.Y)
}

func TestSecondsInRadians(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		time  time.Time
		angle float64
	}{
		{simpleTime(0, 0, 30), math.Pi},
		{simpleTime(0, 0, 0), 0},
		{simpleTime(0, 0, 45), (math.Pi / 2) * 3},
		{simpleTime(0, 0, 7), (math.Pi / 30) * 7},
	} {
		t.Run(testName(c.time), func(t *testing.T) {
			if got := secondsInRadians(c.time); c.angle != got {
				t.Fatalf("want radian %v, got %v", c.angle, got)
			}
		})
	}
}

func TestSecondHandPoint(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 30), Point{0, -1}},
		{simpleTime(0, 0, 45), Point{-1, 0}},
	} {
		t.Run(testName(c.time), func(t *testing.T) {
			if got := secondHandPoint(c.time); !roughlyEqualPoint(got, c.point) {
				t.Fatalf("want point %v, got %v", c.point, got)
			}
		})
	}
}

func TestMinutesInRadians(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		time  time.Time
		angle float64
	}{
		{simpleTime(0, 30, 0), math.Pi},
		{simpleTime(0, 0, 7), 7 * (math.Pi / (30 * 60))},
	} {
		t.Run(testName(c.time), func(t *testing.T) {
			if got := minutesInRadians(c.time); c.angle != got {
				t.Fatalf("want radian %v, got %v", c.angle, got)
			}
		})
	}
}

func TestMinuteHandPoint(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 30, 0), Point{0, -1}},
		{simpleTime(0, 45, 0), Point{-1, 0}},
	} {
		t.Run(testName(c.time), func(t *testing.T) {
			if got := minuteHandPoint(c.time); !roughlyEqualPoint(got, c.point) {
				t.Fatalf("want point %v, got %v", c.point, got)
			}
		})
	}
}

func TestHoursInRadians(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		time  time.Time
		angle float64
	}{
		{simpleTime(6, 0, 0), math.Pi},
		{simpleTime(0, 0, 0), 0},
		{simpleTime(21, 0, 0), math.Pi * 1.5},
		{simpleTime(0, 1, 30), math.Pi / ((6 * 60 * 60) / 90)},
	} {
		t.Run(testName(c.time), func(t *testing.T) {
			if got := hoursInRadians(c.time); !roughlyEqualFloat64(c.angle, got) {
				t.Fatalf("want radian %v, got %v", c.angle, got)
			}
		})
	}
}

func TestHourHandPoint(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(6, 0, 0), Point{0, -1}},
		{simpleTime(21, 0, 0), Point{-1, 0}},
	} {
		t.Run(testName(c.time), func(t *testing.T) {
			if got := hourHandPoint(c.time); !roughlyEqualPoint(c.point, got) {
				t.Fatalf("want point %v, got %v", c.point, got)
			}
		})
	}
}
