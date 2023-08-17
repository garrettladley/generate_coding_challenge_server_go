package domain

import (
	"fmt"
)

type NUID string

func ParseNUID(str string) (*NUID, error) {
	if len(str) != 9 {
		return nil, fmt.Errorf("invalid NUID! Given: %s", str)
	}

	for _, c := range str {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("invalid NUID! Given: %s", str)
		}
	}
	nuid := NUID(str)
	return &nuid, nil
}

func (n NUID) String() string {
	return string(n)
}
