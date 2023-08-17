package domain

func InsertCharAtIndex(str string, char rune, index int) string {
	return str[:index] + string(char) + str[index:]
}
