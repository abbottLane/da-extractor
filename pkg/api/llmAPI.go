package llmapis

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ReadFirstLine reads the first line of a text file
func loadCredentials(textPath string) (string, error) {
	// Open the file
	file, err := os.Open(textPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Scan the first line
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	// If there is no first line, return an empty string and an error
	return "", errors.New("empty file")
}

// Convert a []byte into a json map
func readResponse(resp []byte) (map[string]interface{}, error) {
	var data map[string]interface{}    // Declare a map variable to store the response body
	err := json.Unmarshal(resp, &data) // Decode the response body into the map
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DiscourseCall(messages []map[string]interface{}, credential string) interface{} {
	data := map[string]interface{}{
		"model":    "gpt-4",
		"messages": messages,
	}

	credentials, err := loadCredentials("./llm_key.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+credentials)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	response, err := readResponse(responseBody)
	return response["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"]
}
