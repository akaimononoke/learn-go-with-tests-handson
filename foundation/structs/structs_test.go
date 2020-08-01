package structs

import (
	"testing"
)

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0
	if got != want {
		t.Errorf("got %.2f, want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	t.Parallel()
	for desc, tt := range map[string]struct {
		shape Shape
		want  float64
	}{
		"rectangle": {shape: Rectangle{12, 6}, want: 72.0},
		"circle":    {shape: Circle{10}, want: 314.1592653589793},
		"triangle":  {shape: Triangle{12, 6}, want: 36.0},
	} {
		t.Run(desc, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.want {
				t.Errorf("%#v got %v, want %v", tt.shape, got, tt.want)
			}
		})
	}
}
