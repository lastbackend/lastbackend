package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/scanner"
)

func LoadDockerPatterns(path string) ([]string, error) {

	var excludes []string

	reader, err := os.Open(fmt.Sprintf("%s/.dockerignore", path))
	if err != nil {
		if os.IsNotExist(err) {
			return excludes, nil
		}
		return excludes, nil
	}

	if reader == nil {
		return excludes, nil
	}
	defer reader.Close()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		if pattern == "" {
			continue
		}
		pattern = filepath.Clean(pattern)
		pattern = filepath.ToSlash(pattern)
		excludes = append(excludes, pattern)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading .dockerignore: %v", err)
	}
	return excludes, nil
}

func exclusion(pattern string) bool {
	return pattern[0] == '!'
}

func empty(pattern string) bool {
	return pattern == ""
}

func CleanPatterns(patterns []string) ([]string, [][]string, bool, error) {

	cleanedPatterns := []string{}
	patternDirs := [][]string{}
	exceptions := false
	for _, pattern := range patterns {

		pattern = strings.TrimSpace(pattern)
		if empty(pattern) {
			continue
		}
		if exclusion(pattern) {
			if len(pattern) == 1 {
				return nil, nil, false, errors.New("Illegal exclusion pattern: !")
			}
			exceptions = true
		}
		pattern = filepath.Clean(pattern)
		cleanedPatterns = append(cleanedPatterns, pattern)
		if exclusion(pattern) {
			pattern = pattern[1:]
		}
		patternDirs = append(patternDirs, strings.Split(pattern, string(os.PathSeparator)))
	}

	return cleanedPatterns, patternDirs, exceptions, nil
}

func Matches(file string, patterns []string) (bool, error) {
	file = filepath.Clean(file)

	if file == "." {
		return false, nil
	}

	patterns, patDirs, _, err := CleanPatterns(patterns)
	if err != nil {
		return false, err
	}

	return OptimizedMatches(file, patterns, patDirs)
}

func OptimizedMatches(file string, patterns []string, patDirs [][]string) (bool, error) {
	matched := false
	file = filepath.FromSlash(file)
	parentPath := filepath.Dir(file)
	parentPathDirs := strings.Split(parentPath, string(os.PathSeparator))

	for i, pattern := range patterns {
		negative := false

		if exclusion(pattern) {
			negative = true
			pattern = pattern[1:]
		}

		match, err := regexpMatch(pattern, file)
		if err != nil {
			return false, fmt.Errorf("Error in pattern (%s): %s", pattern, err)
		}

		if !match && parentPath != "." {
			if len(patDirs[i]) <= len(parentPathDirs) {
				match, _ = regexpMatch(strings.Join(patDirs[i], string(os.PathSeparator)),
					strings.Join(parentPathDirs[:len(patDirs[i])], string(os.PathSeparator)))
			}
		}

		if match {
			matched = !negative
		}
	}

	return matched, nil
}

func regexpMatch(pattern, path string) (bool, error) {
	regStr := "^"

	if _, err := filepath.Match(pattern, path); err != nil {
		return false, err
	}

	var scan scanner.Scanner
	scan.Init(strings.NewReader(pattern))

	sl := string(os.PathSeparator)
	escSL := sl
	if sl == `\` {
		escSL += `\`
	}

	for scan.Peek() != scanner.EOF {
		ch := scan.Next()

		if ch == '*' {
			if scan.Peek() == '*' {
				scan.Next()

				if scan.Peek() == scanner.EOF {
					regStr += ".*"
				} else {
					regStr += "((.*" + escSL + ")|([^" + escSL + "]*))"
				}

				if string(scan.Peek()) == sl {
					scan.Next()
				}
			} else {
				regStr += "[^" + escSL + "]*"
			}
		} else if ch == '?' {
			regStr += "[^" + escSL + "]"
		} else if strings.Index(".$", string(ch)) != -1 {
			regStr += `\` + string(ch)
		} else if ch == '\\' {
			if sl == `\` {
				regStr += escSL
				continue
			}
			if scan.Peek() != scanner.EOF {
				regStr += `\` + string(scan.Next())
			} else {
				regStr += `\`
			}
		} else {
			regStr += string(ch)
		}
	}

	regStr += "$"

	res, err := regexp.MatchString(regStr, path)

	if err != nil {
		err = filepath.ErrBadPattern
	}

	return res, err
}
