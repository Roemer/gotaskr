// Package gitlab contains helper methods to work with GitLab.
package gitlab

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/roemer/gotaskr/goext"
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

// GitLabReport defines the data for the quality report for GitLab.
// See https://docs.gitlab.com/ee/ci/testing/code_quality.html#implementing-a-custom-tool for details
type GitLabReport struct {
	Entries []GitLabCodeQualityEntry
}

type GitLabCodeQualityEntry struct {
	Description string                    `json:"description"`
	Fingerprint string                    `json:"fingerprint"`
	Severity    string                    `json:"severity"`
	Location    GitLabCodeQualityLocation `json:"location"`
}

type GitLabCodeQualityLocation struct {
	Path  string                 `json:"path"`
	Lines GitLabCodeQualityLines `json:"lines"`
}

type GitLabCodeQualityLines struct {
	Begin int64 `json:"begin"`
	End   int64 `json:"end"`
}

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

func ConvertEsLintReportToGitLabReport(esLintReport *EsLintReport) (*GitLabReport, error) {
	gitLabReport := &GitLabReport{
		Entries: []GitLabCodeQualityEntry{},
	}
	for _, esLintFile := range esLintReport.Files {

		for _, esLintMessage := range esLintFile.Messages {
			gitlab := GitLabCodeQualityEntry{
				Description: esLintMessage.Message,
				Severity:    goext.Ternary(esLintMessage.Severity == 2, "critical", "info"),
				Location: GitLabCodeQualityLocation{
					Path: esLintFile.FilePath,
					Lines: GitLabCodeQualityLines{
						Begin: esLintMessage.Line,
						End:   esLintMessage.EndLine,
					},
				},
			}
			gitLabReport.Entries = append(gitLabReport.Entries, gitlab)
		}
	}
	return gitLabReport, nil
}

func WriteGitLabReport(gitLabReport *GitLabReport, outputFilePath string) error {
	data, err := json.MarshalIndent(gitLabReport.Entries, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(outputFilePath, data, 0755); err != nil {
		return err
	}
	return nil
}
