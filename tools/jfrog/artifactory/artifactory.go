package artifactory

import (
	"os"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

// ArtifactorySearchResult represents a search result item.
type ArtifactorySearchResult struct {
	Repo       string                `json:"repo"`
	Path       string                `json:"path"`
	Name       string                `json:"name"`
	Created    time.Time             `json:"created"`
	Modified   time.Time             `json:"modified"`
	Type       string                `json:"type"`
	Size       int                   `json:"size"`
	Properties []ArtifactoryProperty `json:"properties"`
}

// ArtifactoryProperty represents a property of a search result item.
type ArtifactoryProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateManager creates a basic manager to interact with artifactory.
func CreateManager(baseUrl string, apiKey string) (artifactory.ArtifactoryServicesManager, error) {
	artifactoryDetails := auth.NewArtifactoryDetails()
	artifactoryDetails.SetUrl(baseUrl)
	artifactoryDetails.SetApiKey(apiKey)

	configBuilder, err := config.NewConfigBuilder().
		SetServiceDetails(artifactoryDetails).
		Build()
	if err != nil {
		return nil, err
	}

	artifactoryManager, err := artifactory.New(configBuilder)
	return artifactoryManager, err
}

// HasSearchResults performs the given search and returns a boolean if the search contains any items or not.
func HasSearchResults(artifactoryManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) (bool, error) {
	reader, err := artifactoryManager.SearchFiles(searchParams)
	if err != nil {
		return false, err
	}
	defer reader.Close()

	return !reader.IsEmpty(), nil
}

// GetSearchResults performs the given search and returns the found items.
func GetSearchResults(artifactoryManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) ([]ArtifactorySearchResult, error) {
	searchResultItems := []ArtifactorySearchResult{}

	reader, err := artifactoryManager.SearchFiles(searchParams)
	if err != nil {
		return searchResultItems, err
	}
	defer reader.Close()

	searchResultItems = GetResultItemsFromReader(reader)

	return searchResultItems, nil
}

// GetResultItemsFromReader is a helper method that converts the search results into typed search result items.
func GetResultItemsFromReader(reader *content.ContentReader) []ArtifactorySearchResult {
	searchResultItems := []ArtifactorySearchResult{}
	if reader != nil {
		for searchResultItem := new(ArtifactorySearchResult); reader.NextRecord(searchResultItem) == nil; searchResultItem = new(ArtifactorySearchResult) {
			searchResultItems = append(searchResultItems, *searchResultItem)
		}
	}
	return searchResultItems
}

// GetContentFromReader is a debug method that returns the raw content of the reader.
func GetContentFromReader(reader *content.ContentReader) string {
	if reader == nil {
		return ""
	}
	content, _ := os.ReadFile(reader.GetFilesPaths()[0])
	return string(content)
}
