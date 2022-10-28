package gttools

import (
	"fmt"
	"os"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"golang.org/x/exp/slices"
)

// ArtifactoryTool provides access to the helper methods for Artifactory.
type ArtifactoryTool struct {
}

// ArtifactorySearchResultItem represents a search result item.
type ArtifactorySearchResultItem struct {
	Repo       string                            `json:"repo"`
	Path       string                            `json:"path"`
	Name       string                            `json:"name"`
	Created    time.Time                         `json:"created"`
	Modified   time.Time                         `json:"modified"`
	Type       string                            `json:"type"`
	Size       int                               `json:"size"`
	Properties []ArtifactorySearchResultProperty `json:"properties"`
}

// ArtifactorySearchResultProperty represents a property of a search result item.
type ArtifactorySearchResultProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetProperty returns the value of the property with the given key and a bool, if the property exists or not.
func (searchResult *ArtifactorySearchResultItem) GetProperty(key string) (string, bool) {
	index := slices.IndexFunc(searchResult.Properties, func(p ArtifactorySearchResultProperty) bool { return p.Key == key })
	if index >= 0 {
		return searchResult.Properties[index].Value, true
	}
	return "", false
}

// HasProperty returns a bool if the property with the given key exists or not.
func (searchResult *ArtifactorySearchResultItem) HasProperty(key string) bool {
	_, hasProperty := searchResult.GetProperty(key)
	return hasProperty
}

// CreateManager creates a basic manager to interact with artifactory.
func (tool *ArtifactoryTool) CreateManager(baseUrl string, apiKey string) (artifactory.ArtifactoryServicesManager, error) {
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
func (tool *ArtifactoryTool) HasSearchResults(artifactoryManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) (bool, error) {
	reader, err := artifactoryManager.SearchFiles(searchParams)
	if err != nil {
		return false, err
	}
	defer reader.Close()

	return !reader.IsEmpty(), nil
}

// GetSearchResults performs the given search and returns the found items.
func (tool *ArtifactoryTool) GetSearchResults(artifactoryManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) ([]*ArtifactorySearchResultItem, error) {
	searchResultItems := []*ArtifactorySearchResultItem{}

	reader, err := artifactoryManager.SearchFiles(searchParams)
	if err != nil {
		return searchResultItems, err
	}
	defer reader.Close()

	searchResultItems = tool.GetResultItemsFromReader(reader)

	return searchResultItems, nil
}

// GetSingleSearchResult returns a single result or nil if none is found and an error, if multiple were found.
func (tool *ArtifactoryTool) GetSingleSearchResult(artifactoryManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) (*ArtifactorySearchResultItem, error) {
	searchResultItems, err := tool.GetSearchResults(artifactoryManager, searchParams)
	if err != nil {
		return nil, err
	}
	if len(searchResultItems) > 1 {
		return nil, fmt.Errorf("got more than one item: %d", len(searchResultItems))
	} else if len(searchResultItems) == 1 {
		return searchResultItems[0], nil
	}
	return nil, nil
}

// GetResultItemsFromReader is a helper method that converts the search results into typed search result items.
func (tool *ArtifactoryTool) GetResultItemsFromReader(reader *content.ContentReader) []*ArtifactorySearchResultItem {
	searchResultItems := []*ArtifactorySearchResultItem{}
	if reader != nil {
		for searchResultItem := new(ArtifactorySearchResultItem); reader.NextRecord(searchResultItem) == nil; searchResultItem = new(ArtifactorySearchResultItem) {
			searchResultItems = append(searchResultItems, searchResultItem)
		}
	}
	return searchResultItems
}

// GetContentFromReader is a debug method that returns the raw content of the reader.
func (tool *ArtifactoryTool) GetContentFromReader(reader *content.ContentReader) string {
	if reader == nil {
		return ""
	}
	content, _ := os.ReadFile(reader.GetFilesPaths()[0])
	return string(content)
}
