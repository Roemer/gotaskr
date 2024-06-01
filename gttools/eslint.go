package gttools

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// EsLintTool provides access to the helper methods for ES Lint.
type EsLintTool struct {
}

func CreateEsLintTool() *EsLintTool {
	return &EsLintTool{}
}

type EsLintReport struct {
	Files []*EsLintFile
}

type EsLintFile struct {
	FilePath string           `json:"filePath"`
	Messages []*EsLintMessage `json:"messages"`
}

type EsLintMessage struct {
	RuleId    string `json:"ruleId"`
	Severity  int64  `json:"severity"`
	Message   string `json:"message"`
	Line      int64  `json:"line"`
	EndLine   int64  `json:"endLine"`
	Column    int64  `json:"column"`
	EndColumn int64  `json:"endColumn"`
}

// ParseEsLintReport parses the given eslint report (JSON).
func (tool *EsLintTool) ParseEsLintReport(esLintReportPath string) (*EsLintReport, error) {
	jsonFile, err := os.Open(esLintReportPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var esLintReportFiles []*EsLintFile
	err = json.Unmarshal(byteValue, &esLintReportFiles)
	if err != nil {
		return nil, err
	}

	esLintReport := &EsLintReport{
		Files: esLintReportFiles,
	}

	return esLintReport, nil
}

func (tool *EsLintTool) SeverityToString(severity int64) string {
	switch severity {
	case 0:
		return "off"
	case 1:
		return "warn"
	case 2:
		return "error"
	default:
		return fmt.Sprintf("unknown (%d)", severity)
	}
}
