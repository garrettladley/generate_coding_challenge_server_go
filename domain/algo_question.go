package domain

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
)

type Challenge struct {
	Challenge []string
	Solution  []string
}

type EditType int

const (
	Insertion EditType = iota
	Deletion
	Substitution
)

func EditTypes() []EditType {
	return []EditType{
		Insertion,
		Deletion,
		Substitution,
	}
}

func GenerateChallenge(nuid NUID, nRandom int, mandatoryCases []string) Challenge {
	rng := rand.New(rand.NewSource(int64(hashCode(nuid.String()))))
	randomCases := make([]string, nRandom)
	colors := Colors()
	editTypes := EditTypes()
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	for i := 0; i < nRandom; i++ {
		color := colors[rng.Intn(len(colors))]
		colorStr, _ := color.String()
		lenColor := len(colorStr)
		randomCount := rng.Intn(lenColor + 1)
		if randomCount == 0 {
			randomCases[i] = colorStr
			continue
		}

		editType := editTypes[rng.Intn(len(editTypes))]

		switch editType {
		case Deletion:
			randomCases[i] = colorStr[randomCount:]
		case Insertion:
			alphabet := "abcdefghijklmnopqrstuvwxyz"
			colorChars := []rune(colorStr)
			randomChars := make([]rune, randomCount)
			for j := 0; j < randomCount; j++ {
				randomChars[j] = rune(alphabet[rng.Intn(len(alphabet))])
			}
			randomIndices := make([]int, randomCount)
			for j := 0; j < randomCount; j++ {
				randomIndices[j] = rng.Intn(len(colorChars) + 1)
			}
			for j := 0; j < randomCount; j++ {
				colorChars = []rune(InsertCharAtIndex(string(colorChars), randomChars[j], randomIndices[j]))
			}
			randomCases[i] = string(colorChars)
		case Substitution:
			colorChars := []rune(colorStr)
			changedIndices := make([]int, randomCount)
			for j := 0; j < randomCount; j++ {
				changedIndices[j] = rng.Intn(lenColor)
			}
			for j := 0; j < randomCount; j++ {
				originalChar := colorChars[changedIndices[j]]
				var newChar rune
				for {
					newChar = rune(alphabet[rng.Intn(len(alphabet))])
					if newChar != originalChar {
						break
					}
				}
				colorChars[changedIndices[j]] = newChar
			}
			randomCases[i] = string(colorChars)
		}
	}

	allCases := append(mandatoryCases, randomCases...)

	var answers []string
	for _, caseColor := range allCases {
		answer, err := oneEditAway(caseColor)
		if err != nil {
			continue
		}
		answerStr, _ := answer.String()
		answers = append(answers, answerStr)
	}

	return Challenge{
		Challenge: allCases,
		Solution:  answers,
	}
}

func oneEditAway(str string) (Color, error) {
	result, err := oneEditAwayOnlyOneAnswer(str, Colors(), func(c Color) string {
		s, _ := c.String()
		return s
	})

	if err != nil {
		return 0, err
	}

	return result, nil
}

func oneEditAwayOnlyOneAnswer[T any](str string, iterable []T, toString func(T) string) (T, error) {
	var defaultItem T

	for _, item := range iterable {
		if nEditsAway(str, toString(item), 1) {
			return item, nil
		}
	}
	return defaultItem, fmt.Errorf("no valid answer found")
}

func nEditsAway(str1, str2 string, n int) bool {
	diffLen := int(math.Abs(float64(len(str1) - len(str2))))
	if diffLen > n {
		return false
	}

	var shorter, longer string
	if len(str1) > len(str2) {
		shorter, longer = str2, str1
	} else {
		shorter, longer = str1, str2
	}

	var shortPointer, longPointer int
	var editCount int

	for shortPointer < len(shorter) && longPointer < len(longer) {
		if shorter[shortPointer] != longer[longPointer] {
			editCount++
			if editCount > n {
				return false
			}
			if len(shorter) == len(longer) {
				shortPointer++
			}
		} else {
			shortPointer++
		}
		longPointer++
	}

	return editCount <= n
}

func hashCode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
