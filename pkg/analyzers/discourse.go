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

func (da *DiscourseAnalyzer) Analyze(text string, tagset []string) string {
	// Define the discourse task directives for the LLM
	instructions := "You are an agent that performs discourse analysis on text. You are given a prompt and you must first split the given text into sentences separated by newlines. For each line you label the discourse function of the sentence."
	formatInstructions := " The format of the response should be JSON with the following keys: 'sentences' and 'discourse'. The 'sentences' key should be a list of strings, each string being a sentence from the input text. The 'discourse' key should be a list of strings, each string being the discourse function of the corresponding sentence. The discourse functions should be one of the following: "
	// loop through the set of strings and add them to a new string newstring
	joinedtagset := ""
	for _, tag := range tagset {
		joinedtagset = joinedtagset + tag + ", "
	}
	// remove the last comma and space from the newstring, add a space to the beginning of the list
	joinedtagset = " " + joinedtagset[:len(joinedtagset)-2]

	// Make a call to the discourse API
	// The pkg/api/llmAPI.go file has a function called discourseCall that makes a call to the discourse API
	messages := []map[string]interface{}{
		{"role": "system",
			"content": instructions + formatInstructions + joinedtagset},
		{"role": "user", "content": text},
	}
	result := llm.DiscourseCall(messages, da.apiKey)
	// return the resul"t of the discourse call
	return result.(string)
}
