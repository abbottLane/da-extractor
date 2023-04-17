package discourse

import (
	llm "daextractor/pkg/api"
)

// make a class with methods for the discourse analyzer
type DiscourseAnalyzer struct {
	apiService string
	apiKey     string
}

func NewDiscourseAnalyzer(apiService string) *DiscourseAnalyzer {
	// print out the apiService and apiKey
	return &DiscourseAnalyzer{apiService: apiService}
}

func (da *DiscourseAnalyzer) Analyze(text string) string {
	// Make a call to the discourse API
	// The pkg/api/llmAPI.go file has a function called discourseCall that makes a call to the discourse API

	result := llm.DiscourseCall(text, da.apiKey)
	// return the resul"t of the discourse call
	return result.(string)
}
