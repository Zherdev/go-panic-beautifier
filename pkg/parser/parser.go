package parser

import (
	"fmt"
	"strings"
)

type lineType int

const (
	headerLine = lineType(iota)
	signalLine
	goroutineLine
	pathLine
	unkLine
)

type PanicLogParser struct {
	src string
}

func NewPanicLogParser(src string) *PanicLogParser {
	return &PanicLogParser{src: src}
}

func (p *PanicLogParser) Parse() string {
	p.preprocess()
	p.parseLines()
	p.parseSpecSymbols()

	return p.src
}

func (p *PanicLogParser) preprocess() {
	p.src = strings.ReplaceAll(p.src, "\n", "\\n")
	p.src = strings.ReplaceAll(p.src, "\t", "\\t")
	p.src = strings.ReplaceAll(p.src, "    ", "\\t")
}

func (p *PanicLogParser) parseSpecSymbols() {
	p.src = strings.ReplaceAll(p.src, "\\n", "<br>")
	p.src = strings.ReplaceAll(p.src, "\\t", "&emsp;")
}

func (p *PanicLogParser) parseLines() {
	lines := strings.Split(p.src, "\\n")
	for idx, line := range lines {
		lineType := p.getLineType(line)
		lines[idx] = p.processLine(line, lineType)
		fmt.Printf("type %d", lineType)
	}

	p.src = strings.Join(lines, "\\n")
}

func (p *PanicLogParser) getLineType(line string) lineType {
	line = strings.TrimSpace(line)

	if !strings.ContainsRune(line, '/')  && !strings.ContainsRune(line, '[') &&
		strings.ContainsRune(line, ':') {
		return headerLine
	}

	if !strings.ContainsRune(line, '/') && strings.HasPrefix(line, "[signal ") &&
		strings.ContainsRune(line, ']') {
		return signalLine
	}

	if strings.HasPrefix(line, "goroutine") {
		return goroutineLine
	}

	if strings.HasPrefix(line, "/") {
		return pathLine
	}

	return unkLine
}

func (p *PanicLogParser) processLine(line string, t lineType) string {
	switch t {
	case headerLine:
		// highlight error message,
		// ex: "invalid memory address or nil pointer dereference"

		parts := strings.Split(line, ":")

		partIdx := len(parts) - 1
		for partIdx > 0 && len(strings.TrimSpace(parts[partIdx])) == 0 {
			partIdx--
		}

		parts[partIdx] = "<span class=\"header-error-msg\">" + parts[partIdx] + "</span>"
		return strings.Join(parts, ":")

	case signalLine:
		// highlight signal message,
		// ex: "SIGSEGV: segmentation violation"
		sigStartPos := len("[signal ")
		sigEndPos := strings.Index(line, "code")
		if sigEndPos == -1 {
			sigEndPos = strings.Index(line, "]") // != -1
		}

		prefix := line[:sigStartPos]
		suffix := line[sigEndPos:]
		sig := line[sigStartPos:sigEndPos]

		return prefix + "<span class=\"signal-msg\">" + sig + "</span>" + suffix

	case goroutineLine:
		// highlight full goroutine line,
		// ex: "goroutine 1 [running]:"
		return "<span class=\"goroutine-line\">" + line + "</span>"

	case pathLine:
		// highlight line number, ex: "...:126"
		lineNumStartPos := strings.LastIndex(line, ":")
		if lineNumStartPos != -1 && lineNumStartPos < len(strings.TrimSpace(line)) - 1 {
			line = fmt.Sprintf("%s<span class=\"line-num\">%s</span>",
				line[:lineNumStartPos+1], line[lineNumStartPos+1:])
		}

		// highlight filename, ex: ".../file.go:..."
		fileStartPos := strings.LastIndex(line, "/")
		fileEndPos := len(line)
		if lineNumStartPos != -1 {
			fileEndPos = lineNumStartPos + 1
		}
		if fileStartPos != -1 && fileStartPos < len(strings.TrimSpace(line)) - 1 {
			line = line[:fileStartPos+1] + "<span class=\"filename\">" + line[lineNumStartPos+1:] + "</span>" + line[:fileEndPos]
		}

		// highlight path, ex: "/Users/i.zherdev/..."
		return line

	default:
		return line
	}
}
