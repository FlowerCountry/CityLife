package locale

import (
	"bufio"
	"strings"
)

// parseINI 解析INI格式内容
func (l *Locale) parseINI(content string) error {
	scanner := bufio.NewScanner(strings.NewReader(content))
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || line[0] == '#' || line[0] == ';' {
			continue
		}

		// 解析 section [section_name]
		if line[0] == '[' {
			closePos := strings.Index(line, "]")
			if closePos != -1 {
				currentSection = line[1:closePos]
				// 确保 section 存在
				if _, ok := l.sections[currentSection]; !ok {
					l.sections[currentSection] = make(map[string]string)
				}
			}
			continue
		}

		// 解析 key=value
		eqPos := strings.Index(line, "=")
		if eqPos == -1 {
			continue
		}

		key := strings.TrimSpace(line[:eqPos])
		value := strings.TrimSpace(line[eqPos+1:])

		// 存储到当前 section
		if currentSection != "" && key != "" {
			l.sections[currentSection][key] = value
		}
	}

	return scanner.Err()
}
