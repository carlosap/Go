package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const ZULU_PARSE_FORMAT = "2006-01-02T15:04:05Z"
const ZULU_PARSE_FORMAT_MS = "2006-01-02T15:04:05.000Z"
const LONG_PARSE_FORMAT = "2006-01-02T15:04:05.999999Z"
const TIMEZONE_PARSE_FORMAT = "2006-01-02T15:04:05-07:00"
const DAYS_FROM_1900_TO_1970 = 25569.0
const SECONDS_PER_DAY = 86400.0
const NS_PER_MS = 1000000

var groupingReg = regexp.MustCompile(`(\D+)(\d+)`)
var mtxReg = regexp.MustCompile(`MTX.*`)
var multiWhiteReg = regexp.MustCompile(`[\s]{2,}`)
var punctReg = regexp.MustCompile(`[!"#$%&'()*+,-.:;<=>?@\[\]^_\{|}~]`)
var whiteReg = regexp.MustCompile(`[\s]`)

func ConvertFromDaysToNanoSeconds(timeDays float64) int64 {
	daysSince1970 := timeDays - DAYS_FROM_1900_TO_1970
	secondsSince1970 := daysSince1970 * SECONDS_PER_DAY
	return int64(secondsSince1970)
}

func FloatTime(t int64) float64 {
	return float64(t)
}

func GetSource(uid string) string {
	if lastIdx := strings.LastIndex(uid, "."); lastIdx > 0 {
		lastPart := uid[lastIdx+1:]
		if strings.Contains(uid, "mPT.") || strings.Contains(uid, "SCG.") || MatchesShortCallsign(lastPart) {
			return "Video"
		}
		return uid[:lastIdx]
	}

	if strings.Contains(uid, "mPT.") || strings.Contains(uid, "SCG.") || MatchesShortCallsign(uid) {
		return "Video"
	}

	return uid
}

// Remove punctuation and duplicate whitespace
func punctCleanWhite(callsign string) string {
	callsign = punctReg.ReplaceAllLiteralString(callsign, "")
	callsign = multiWhiteReg.ReplaceAllLiteralString(callsign, " ")
	return strings.TrimSpace(callsign)
}

func CurrentTime() int64 {
	return time.Now().UTC().Unix()
}

func RemoveSPI(uid string) string {
	lowerUid := strings.ToLower(uid)
	lastSpi := strings.LastIndex(lowerUid, "spi")
	if lastSpi > 1 {
		return uid[:lastSpi-1]
	}
	return uid
}

func IsDigit(part byte) bool {
	return '0' <= part && part <= '9'
}

func MatchesShortCallsign(callsign string) bool {
	if len(callsign) == 4 {
		return !IsDigit(callsign[0]) && !IsDigit(callsign[1]) && IsDigit(callsign[2]) && IsDigit(callsign[3])
	}
	return false
}

func IsSub(cotType string) bool {
	if strings.HasPrefix(cotType, "a") {
		if len(cotType) >= 5 {
			return cotType[4] == 'U'
		}
	}
	return false
}

func IsSea(cotType string) bool {
	if strings.HasPrefix(cotType, "a") {
		if len(cotType) >= 5 {
			return cotType[4] == 'S'
		}
	}
	return false
}

func IsGround(cotType string) bool {
	if strings.HasPrefix(cotType, "a") {
		if len(cotType) >= 5 {
			return cotType[4] == 'G'
		}
	}
	return false
}

func IsUnknown(cotType string) bool {
	if strings.HasPrefix(cotType, "a") {
		if len(cotType) >= 5 {
			return cotType[4] == 'Z'
		}
	}
	return false
}

func IsSpace(cotType string) bool {
	if strings.HasPrefix(cotType, "a") {
		if len(cotType) >= 5 {
			return cotType[4] == 'P'
		}
	}
	return false
}

func IsAir(cotType string) bool {
	if strings.HasPrefix(cotType, "a") {
		if len(cotType) >= 5 {
			return cotType[4] == 'A'
		}
	}
	return false
}

func IsSPI(varType string) bool {
	return strings.HasSuffix(varType, "s-p-i")
}

func ConvertToTimezoneString(Unix int64) string {
	timeValue := time.Unix(Unix, 0)
	return timeValue.Format(TIMEZONE_PARSE_FORMAT)
}

func ConvertToZuluString(Unix int64) string {
	timeValue := time.Unix(Unix, 0)
	return timeValue.Format(LONG_PARSE_FORMAT)
}

//ParseTime parses a CoT formatted time string into the Unix representation
func ParseTime(timeString string) int64 {
	parsedUnixTime := int64(-1)
	var parsedTimeStruct time.Time
	var err error

	if strings.HasSuffix(timeString, "Z") {
		parsedTimeStruct, err = time.Parse(ZULU_PARSE_FORMAT, timeString)
	} else {
		parsedTimeStruct, err = time.Parse(TIMEZONE_PARSE_FORMAT, timeString)
	}

	if err != nil {
		fmt.Printf("Parse error for time %v : %v", timeString, err)
	} else {
		parsedUnixTime = parsedTimeStruct.Unix()
	}
	return parsedUnixTime
}

// Test if one callsign matches another
func IsMatchingCallsign(callsign1, callsign2 string) bool {
	norm1, norm2 := NormalizeCallsign(callsign1), NormalizeCallsign(callsign2)
	return norm1 != "" && norm2 != "" && norm1 == norm2
}

func GetShortCallsign(callsign string) string {
	// KILLER 11 = KR11
	normCallsign := ""
	fixedCallsign := NormalizeCallsign(callsign)
	callsignParts := strings.Split(fixedCallsign, " ")
	log.Printf("%s: CSP: %+v", callsign, callsignParts)
	if len(callsignParts) >= 2 {
		wordPart := callsignParts[0]
		numberPart := callsignParts[1]
		if len(wordPart) >= 2 && len(numberPart) == 2 && IsDigit(numberPart[0]) && IsDigit(numberPart[1]) {
			normCallsign = fmt.Sprintf("%s%s%s", string(wordPart[0]), string(wordPart[len(wordPart)-1]), numberPart)
		}
	}
	return normCallsign
}

func IsMatchingTadil(tadil1, tadil2 string) bool {
	if tadil1 == "" || tadil2 == "" {
		return false
	}

	return punctCleanWhite(tadil1) == punctCleanWhite(tadil2)
}

// Check if two CoT types are compatible. Special case spi
func IsCompatibleType(var1Type, var2Type string) bool {
	if var1Type == "" || var2Type == "" {
		return true
	}

	if var1Type == var2Type {
		return true
	} else if strings.HasPrefix(var1Type, "a") && strings.HasPrefix(var2Type, "a") {
		// TODO: HACK: Currently match ground with air
		return true
		// Check if we're dealwith with a SPI or a normal platform
		//		if len(var1Type) >= 5 && len(var2Type) >= 5 {
		//			return var1Type[4] == var2Type[4]
		//		}
	} else if IsSPI(var1Type) && IsSPI(var2Type) {
		return true
	}

	return false
}

func IsRadarTrack(uid string) bool {
	uidParts := strings.Split(uid, ".")
	uidLen := len(uidParts)
	if uidLen == 3 && len(uidParts[1]) == 1 && len(uidParts[2]) == 5 {
		return true
	}
	return false
}

func ChooseCanon(name1, name2 string) string {
	return Canonicalize([]string{name1, name2})
}

func Canonicalize(nameList []string) string {
	canon := ""
	sort.Strings(nameList)
	for _, name := range nameList {
		if len(name) > len(canon) {
			canon = name
		}
	}
	return canon
}

func MultiValueDedupe(nameList []string) string {
	sort.Strings(nameList)
	canonValue := Canonicalize(nameList)
	for _, name := range nameList {
		name = punctCleanWhite(name)
		name = whiteReg.ReplaceAllLiteralString(name, "")
		name = groupingReg.ReplaceAllString(name, "$1 $2")
		name = strings.TrimSpace(name)
	}

	return canonValue
}

func ConvertTypeToEnum(cotType string) string {
	enumType := ""
	if len(cotType) >= 5 {
		if strings.HasPrefix(cotType, "b-m-p-s-p-i") {
			enumType = "SPI"
		}
		if cotType[4] == 'a' {
			enumType = "Aircraft"
		} else if cotType[4] == 'g' {
			enumType = "Ground"
		}
		// TODO: Convert CotEnum to type for icons
	}

	return enumType
}

func Convert2525ToCotType(MilSym string) string {
	if MilSym == "" || len(MilSym) < 3 {
		return "a-u-Z"
	}
	// Descriptor
	descriptor := strings.ToLower(string(MilSym[1]))
	battleDimension := string(MilSym[2])
	// SNAP
	symbol := fmt.Sprintf("a-%s-%s", descriptor, battleDimension)
	if len(MilSym) > 9 {
		functionId := strings.TrimRight(MilSym[4:10], "-")
		splitFunIdList := strings.Split(functionId, "")
		if len(splitFunIdList) > 0 {
			restPart := strings.Join(splitFunIdList, "-")
			symbol = fmt.Sprintf("%s-%s", symbol, restPart)
		}
	}
	return symbol
}

func ConvertCotTypeTo2525(cotType string) string {
	if cotType == "" || len(cotType) < 5 {
		return "SUZP"
	}
	parts := strings.Split(cotType, "-")
	// a-f-A-M-F-Q = SFAPMFQ
	joined := strings.Join(parts, "")
	uncap := fmt.Sprintf("S%sP%s", joined[1:3], joined[3:])
	return strings.ToUpper(uncap)
}

func GetCallsignPrefix(callsign string) string {
	normCallsign := NormalizeCallsign(callsign)
	if normCallsign != "" {
		groupedCallsign := groupingReg.ReplaceAllString(normCallsign, "$1")
		groupedCallsign = strings.TrimSpace(groupedCallsign)
		return groupedCallsign
	}
	return ""
}

// Normalize a callsign to have the format: "[String] [Digits]"
func NormalizeCallsign(callsign string) string {
	if callsign == "" {
		return callsign
	}

	callsign = strings.ToUpper(callsign)
	callsign = strings.TrimSpace(callsign)

	_, err := strconv.Atoi(callsign)
	if err == nil {
		return callsign
	}

	if strings.HasPrefix(callsign, "ISAF") || strings.HasPrefix(callsign, "ISF") {
		return callsign
	}

	callsign = punctCleanWhite(callsign)
	callsign = mtxReg.ReplaceAllLiteralString(callsign, "")
	callsign = whiteReg.ReplaceAllLiteralString(callsign, "")
	callsign = groupingReg.ReplaceAllString(callsign, "$1 $2")
	callsign = strings.TrimSpace(callsign)

	return callsign
}

func ZipData(files map[string]io.Reader) io.ReadSeeker {

	buf := new(bytes.Buffer)
	zipper := zip.NewWriter(buf)

	for k, v := range files {
		z, err := zipper.Create(k)
		if err != nil {
			log.Printf("Error creating zip file: %+v", err)
			return nil
		}
		_, err = io.Copy(z, v)
		if err != nil {
			log.Printf("Error writing zip file: %+v", err)
			return nil
		}
	}
	zipper.Close()
	return bytes.NewReader(buf.Bytes())
}
