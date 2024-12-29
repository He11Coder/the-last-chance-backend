package swearWordsDetector

import (
	"bufio"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
)

var re *regexp.Regexp

func uploadSwearsList() ([]string, error) {
	var swearWordsList []string

	english_file, err := os.Open("assets/swears/english_swears.txt")
	if err != nil {
		return nil, err
	}
	defer english_file.Close()

	scanner := bufio.NewScanner(english_file)
	for scanner.Scan() {
		swearWordsList = append(swearWordsList, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	russian_file, err := os.Open("assets/swears/russian_swears.txt")
	if err != nil {
		return nil, err
	}
	defer russian_file.Close()

	scanner = bufio.NewScanner(russian_file)
	for scanner.Scan() {
		swearWordsList = append(swearWordsList, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return swearWordsList, nil
}

func buildRegexpPattern(swearWordsList []string) string {
	var variations []string
	for _, word := range swearWordsList {
		variation := ""
		for _, char := range word {
			escapedChar := regexp.QuoteMeta(string(char))
			variation += "[" + escapedChar
			subs := getSubstitution(char)
			for _, sub := range subs {
				variation += string(sub)
			}
			variation += "][ .\\-]*"
		}
		variation = strings.TrimSuffix(variation, "[ .\\-]*")
		variations = append(variations, variation)
	}

	return strings.Join(variations, "|")
}

func getSubstitution(char rune) []rune {
	switch char {
	case 'a', 'A', 'а', 'А':
		return []rune{'a', 'а', '@'}
	case 'o', 'O', 'о', 'О':
		return []rune{'o', 'о', '0'}
	case 'e', 'E', 'е', 'Е':
		return []rune{'e', 'е', 'ё', '3', 'з'}
	case 'i', 'I':
		return []rune{'i', '1', '!'}
	case 'w', 'W', 'ш', 'Ш':
		return []rune{'w', 'ш'}
	case 't', 'T', 'т', 'Т':
		return []rune{'t', 'т', 'm'}
	case 'y', 'Y', 'у', 'У':
		return []rune{'y', 'у'}
	case 'p', 'P', 'р', 'Р':
		return []rune{'p', 'р'}
	case 's', 'S':
		return []rune{'s', '5', '$'}
	case 'h', 'H', 'н', 'Н':
		return []rune{'h', 'н'}
	case 'k', 'K', 'к', 'К':
		return []rune{'k', 'к'}
	case 'l', 'L':
		return []rune{'l', '1', '!'}
	case 'x', 'X', 'х', 'Х':
		return []rune{'x', 'х', '×', '*'}
	case 'c', 'C', 'с', 'С':
		return []rune{'c', 'с', '('}
	case 'b', 'B', 'в', 'В':
		return []rune{'b', 'в', '8'}
	case 'n', 'N', 'п', 'П':
		return []rune{'n', 'п'}
	case 'm':
		return []rune{'m', 'т'}
	case 'з', 'З', '3':
		return []rune{'з', '3'}
	case 'м', 'М', 'M':
		return []rune{'м', 'm'}
	case 'u', 'U', 'и', 'И':
		return []rune{'u', 'и'}
	default:
		return []rune{char}
	}
}

func BuildAndCompileRegexp() error {
	word_list, err := uploadSwearsList()
	if err != nil {
		return err
	}

	pattern := `(?i)(^|[\P{L}])(` + buildRegexpPattern(word_list) + `)($|[\P{L}])`
	re = regexp.MustCompile(pattern)

	return nil
}

func ContainsSwearWords(input string) bool {
	return re.MatchString(input)
}

func DetectInMultipleInputs(inputs ...string) bool {
	results := make([]bool, len(inputs))
	wg := &sync.WaitGroup{}

	for job_number, input := range inputs {
		wg.Add(1)
		go func(i int, in string) {
			defer wg.Done()

			contains := ContainsSwearWords(in)
			results[i] = contains
		}(job_number, input)
	}

	wg.Wait()

	return slices.Contains(results, true)
}

/*func main() {
	err := BuildAndCompileRegexp()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(DetectInMultipleInputs("привет", "хй", "fu"))
}*/
