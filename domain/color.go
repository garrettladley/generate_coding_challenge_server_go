package domain

import "fmt"

type Color int

const (
	Red Color = iota
	Orange
	Yellow
	Green
	Blue
	Violet
)

func (c Color) String() (string, error) {
	switch c {
	case Red:
		return "red", nil
	case Orange:
		return "orange", nil
	case Yellow:
		return "yellow", nil
	case Green:
		return "green", nil
	case Blue:
		return "blue", nil
	case Violet:
		return "violet", nil
	default:
		return "", fmt.Errorf("invalid color: %d", c)
	}
}

func ParseColor(s string) (Color, error) {
	switch s {
	case "red":
		return Red, nil
	case "orange":
		return Orange, nil
	case "yellow":
		return Yellow, nil
	case "green":
		return Green, nil
	case "blue":
		return Blue, nil
	case "violet":
		return Violet, nil
	default:
		return 0, fmt.Errorf("invalid color: %s", s)
	}
}

func Colors() []Color {
	return []Color{
		Red,
		Orange,
		Yellow,
		Green,
		Blue,
		Violet,
	}
}
