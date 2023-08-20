package domain

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
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

func GenerateChallenge(nRandom int, mandatoryCases []string) Challenge {
	randomCases := make([]string, nRandom)
	colors := Colors()
	editTypes := EditTypes()
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	for i := 0; i < nRandom; i++ {
		randColorIdx, _ := crand.Int(crand.Reader, big.NewInt(int64(len(colors))))
		color := colors[randColorIdx.Int64()]
		colorStr, _ := color.String()
		lenColor := len(colorStr)
		randomCount, _ := crand.Int(crand.Reader, big.NewInt(int64(lenColor+1)))
		if randomCount.Int64() == 0 {
			randomCases[i] = colorStr
			continue
		}

		randEditTypeIdx, _ := crand.Int(crand.Reader, big.NewInt(int64(len(editTypes))))
		editType := editTypes[randEditTypeIdx.Int64()]

		switch editType {
		case Deletion:
			randomCases[i] = colorStr[randomCount.Int64():]
		case Insertion:
			colorChars := []rune(colorStr)
			randomChars := make([]rune, randomCount.Int64())
			for j := 0; j < int(randomCount.Int64()); j++ {
				randomCharIdx, _ := crand.Int(crand.Reader, big.NewInt(int64(len(alphabet))))
				randomChars[j] = rune(alphabet[randomCharIdx.Int64()])
			}
			randomIndices := make([]int, randomCount.Int64())
			for j := 0; j < int(randomCount.Int64()); j++ {
				x, _ := crand.Int(crand.Reader, big.NewInt(int64(len(colorChars))))
				randomIndices[j] = int(x.Int64())
			}
			for j := 0; j < int(randomCount.Int64()); j++ {
				colorChars = []rune(InsertCharAtIndex(string(colorChars), randomChars[j], randomIndices[j]))
			}
			randomCases[i] = string(colorChars)
		case Substitution:
			colorChars := []rune(colorStr)
			changedIndices := make([]int, randomCount.Int64())
			for j := 0; j < int(randomCount.Int64()); j++ {
				randColorIdx, _ := crand.Int(crand.Reader, big.NewInt(int64(lenColor)))
				changedIndices[j] = int(randColorIdx.Int64())
			}
			for j := 0; j < int(randomCount.Int64()); j++ {
				originalChar := colorChars[changedIndices[j]]
				var newChar rune
				for {
					randCharIdx, _ := crand.Int(crand.Reader, big.NewInt(int64(len(alphabet))))
					newChar = rune(alphabet[randCharIdx.Int64()])
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
	mrand.Shuffle(len(allCases), func(i, j int) {
		allCases[i], allCases[j] = allCases[j], allCases[i]
	})

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
