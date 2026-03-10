package frpc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type NumberPair struct {
	First  int64
	Second int64
}

func ParseRangeNumbers(rangeStr string) ([]int64, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	numbers := make([]int64, 0)

	numRanges := strings.Split(rangeStr, ",")
	for _, numRangeStr := range numRanges {
		numArray := strings.Split(numRangeStr, "-")
		switch len(numArray) {
		case 1:
			singleNum, err := strconv.ParseInt(strings.TrimSpace(numArray[0]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("range number is invalid: %v", err)
			}
			numbers = append(numbers, singleNum)
		case 2:
			minValue, err := strconv.ParseInt(strings.TrimSpace(numArray[0]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("range number is invalid: %v", err)
			}
			maxValue, err := strconv.ParseInt(strings.TrimSpace(numArray[1]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("range number is invalid: %v", err)
			}
			if maxValue < minValue {
				return nil, fmt.Errorf("range number is invalid: maxValue < minValue")
			}
			for i := minValue; i <= maxValue; i++ {
				numbers = append(numbers, i)
			}
		default:
			return nil, fmt.Errorf("range number is invalid")
		}
	}

	return numbers, nil
}

func ParseNumberRangePair(firstRangeStr, secondRangeStr string) ([]NumberPair, error) {
	firstRangeNumbers, err := ParseRangeNumbers(firstRangeStr)
	if err != nil {
		return nil, err
	}
	secondRangeNumbers, err := ParseRangeNumbers(secondRangeStr)
	if err != nil {
		return nil, err
	}
	if len(firstRangeNumbers) != len(secondRangeNumbers) {
		return nil, fmt.Errorf("first and second range numbers are not in pairs")
	}

	pairs := make([]NumberPair, 0, len(firstRangeNumbers))
	for i := range firstRangeNumbers {
		pairs = append(pairs, NumberPair{
			First:  firstRangeNumbers[i],
			Second: secondRangeNumbers[i],
		})
	}
	return pairs, nil
}

func FormatPortRange(numbers []int64) string {
	if len(numbers) == 0 {
		return ""
	}

	min := numbers[0]
	max := numbers[0]
	for _, n := range numbers {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	if min == max {
		return fmt.Sprintf("%d", min)
	}
	return fmt.Sprintf("%d-%d", min, max)
}

func ContainsGoTemplate(content string) bool {
	return strings.Contains(content, "{{") && strings.Contains(content, "}}")
}

type TemplateDisplayInfo struct {
	NamePattern string
	Protocol    string
	LocalPorts  string
	RemotePorts string
	RawContent  string
}

func ParseGoTemplateBlock(content string) (*TemplateDisplayInfo, error) {
	info := &TemplateDisplayInfo{
		RawContent: content,
	}

	re := regexp.MustCompile(`parseNumberRangePair\s+"([^"]+)"\s+"([^"]+)"`)
	matches := re.FindAllStringSubmatch(content, -1)
	if len(matches) >= 1 {
		firstRange := strings.Trim(matches[0][1], `"`)
		secondRange := strings.Trim(matches[0][2], `"`)

		firstNumbers, err := ParseRangeNumbers(firstRange)
		if err != nil {
			return nil, err
		}
		secondNumbers, err := ParseRangeNumbers(secondRange)
		if err != nil {
			return nil, err
		}

		info.LocalPorts = FormatPortRange(firstNumbers)
		info.RemotePorts = FormatPortRange(secondNumbers)
	}

	nameRe := regexp.MustCompile(`(?m)name\s*=\s*"([^"]+)"`)
	if match := nameRe.FindStringSubmatch(content); match != nil {
		nameTemplate := match[1]
		nameTemplate = strings.ReplaceAll(nameTemplate, `{{ $v.First }}`, "*")
		nameTemplate = strings.ReplaceAll(nameTemplate, `{{ $v.Second }}`, "*")
		nameTemplate = regexp.MustCompile(`\{\{[^}]+\}\}`).ReplaceAllString(nameTemplate, "*")
		info.NamePattern = strings.Trim(nameTemplate, `"`)
	}

	typeRe := regexp.MustCompile(`(?m)type\s*=\s*"([^"]+)"`)
	if match := typeRe.FindStringSubmatch(content); match != nil {
		info.Protocol = strings.ToUpper(strings.Trim(match[1], `"`))
	}

	if info.NamePattern == "" {
		info.NamePattern = "模板服务"
	}
	if info.Protocol == "" {
		info.Protocol = "TCP"
	}

	return info, nil
}
