// Package eslint contains helper methods to work with ES Lint.
package eslint

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type EsLintReport struct {
	Files []EsLintFile
}

type EsLintFile struct {
	FilePath string          `json:"filePath"`
	Messages []EsLintMessage `json:"messages"`
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

// ParseEsLintReport parses the given eslint report (json).
func ParseEsLintReport(esLintReportPath string) (*EsLintReport, error) {
	jsonFile, err := os.Open(esLintReportPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var esLintReportFiles []EsLintFile
	err = json.Unmarshal(byteValue, &esLintReportFiles)
	if err != nil {
		return nil, err
	}

	esLintReport := &EsLintReport{
		Files: esLintReportFiles,
	}

	return esLintReport, nil
}
