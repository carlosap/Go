package strings

import "strings"

func GetLastValueFromSeparator(value string, separator string) string {
	var retVal string
	if strings.Contains(value, separator) {
		pArray := strings.Split(value, separator)
		retVal = pArray[len(pArray)-1]
	}
	return retVal
}
