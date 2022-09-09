// Package gitlab contains helper methods to work with GitLab.
package gitlab

import (
	"github.com/roemer/gotaskr/goext"
	"github.com/roemer/gotaskr/tools/eslint"
)

// IsRunningOnGitLab returns a flag, if we are currenlty running on gitlab
func IsRunningOnGitLab() bool {
	return goext.EnvExists("GITLAB_CI")
}

// GitLabReport defines the data for the quality report for GitLab.
// See https://docs.gitlab.com/ee/ci/testing/code_quality.html#implementing-a-custom-tool for details
type GitLabReport struct {
	Entries []*GitLabCodeQualityEntry
}

type GitLabCodeQualityEntry struct {
	Description string                     `json:"description"`
	Fingerprint string                     `json:"fingerprint"`
	Severity    string                     `json:"severity"`
	Location    *GitLabCodeQualityLocation `json:"location"`
}

type GitLabCodeQualityLocation struct {
	Path  string                  `json:"path"`
	Lines *GitLabCodeQualityLines `json:"lines"`
}

type GitLabCodeQualityLines struct {
	Begin int64 `json:"begin"`
	End   int64 `json:"end"`
}

// ConvertEsLintReportToGitLabReport converts the given eslint.EsLintReport to a GitLabReport.
func ConvertEsLintReportToGitLabReport(esLintReport *eslint.EsLintReport) *GitLabReport {
	gitLabReport := &GitLabReport{
		Entries: []*GitLabCodeQualityEntry{},
	}
	for _, esLintFile := range esLintReport.Files {

		for _, esLintMessage := range esLintFile.Messages {
			gitlab := &GitLabCodeQualityEntry{
				Description: esLintMessage.Message,
				Severity:    goext.Ternary(esLintMessage.Severity == 2, "critical", "info"),
				Location: &GitLabCodeQualityLocation{
					Path: esLintFile.FilePath,
					Lines: &GitLabCodeQualityLines{
						Begin: esLintMessage.Line,
						End:   esLintMessage.EndLine,
					},
				},
			}
			gitLabReport.Entries = append(gitLabReport.Entries, gitlab)
		}
	}
	return gitLabReport
}

// MergeGitLabReports merges the given GitLabReports to a single GitLabReport.
func MergeGitLabReports(gitLabReports []*GitLabReport) *GitLabReport {
	gitLabReport := &GitLabReport{
		Entries: []*GitLabCodeQualityEntry{},
	}
	for _, r := range gitLabReports {
		gitLabReport.Entries = append(gitLabReport.Entries, r.Entries...)
	}
	return gitLabReport
}

// WriteGitLabReport writes the GitLabReport into a json file.
func WriteGitLabReport(gitLabReport *GitLabReport, outputFilePath string) error {
	return goext.WriteJsonToFile(gitLabReport.Entries, outputFilePath, true)
}
