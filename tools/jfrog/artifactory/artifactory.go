package artifactory

import (
	"os"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"golang.org/x/exp/slices"
)

// SearchResultItem represents a search result item.
type SearchResultItem struct {
	Repo       string                 `json:"repo"`
	Path       string                 `json:"path"`
	Name       string                 `json:"name"`
	Created    time.Time              `json:"created"`
	Modified   time.Time              `json:"modified"`
	Type       string                 `json:"type"`
	Size       int                    `json:"size"`
	Properties []SearchResultProperty `json:"properties"`
}

// SearchResultProperty represents a property of a search result item.
type SearchResultProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetProperty returns the value of the property with the given key and a bool, if the property exists or not.
func (searchResult *SearchResultItem) GetProperty(key string) (string, bool) {
	index := slices.IndexFunc(searchResult.Properties, func(p SearchResultProperty) bool { return p.Key == key })
	if index >= 0 {
		return searchResult.Properties[index].Value, true
	}
	return "", false
}

// HasProperty returns a bool if the property with the given key exists or not.
func (searchResult *SearchResultItem) HasProperty(key string) bool {
	_, hasProperty := searchResult.GetProperty(key)
	return hasProperty
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

// HasSearchResults performs the given search and returns a bool if the search contains any items or not.
func HasSearchResults(artifactoryManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) (bool, error) {
	reader, err := artifactoryManager.SearchFiles(searchParams)
	if err != nil {
		return false, err
	}
	defer reader.Close()

	return !reader.IsEmpty(), nil
}

// GetSearchResults performs the given search and returns the found items.
func GetSearchResults(artifactoryManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) ([]SearchResultItem, error) {
	searchResultItems := []SearchResultItem{}

	reader, err := artifactoryManager.SearchFiles(searchParams)
	if err != nil {
		return searchResultItems, err
	}
	defer reader.Close()

	searchResultItems = GetResultItemsFromReader(reader)

	return searchResultItems, nil
}

// GetResultItemsFromReader is a helper method that converts the search results into typed search result items.
func GetResultItemsFromReader(reader *content.ContentReader) []SearchResultItem {
	searchResultItems := []SearchResultItem{}
	if reader != nil {
		for searchResultItem := new(SearchResultItem); reader.NextRecord(searchResultItem) == nil; searchResultItem = new(SearchResultItem) {
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
