package result

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// On windows the pty package is failing to obtain stderr text when a test fails,
// so this is our workaround: when a test fails, it writes the result to a file,
// and then the test runner reads the file and checks the result

const PathEnvVar = "LAZYGIT_INTEGRATION_TEST_RESULT_PATH"

type IntegrationTestResult struct {
	Success bool
	Message string
}

func LogFailure(message string) error {
	return writeResult(failure(message))
}

func LogSuccess() error {
	return writeResult(success())
}

func failure(message string) IntegrationTestResult {
	return IntegrationTestResult{Success: false, Message: message}
}

func success() IntegrationTestResult {
	return IntegrationTestResult{Success: true}
}

func writeResult(result IntegrationTestResult) error {
	resultPath := os.Getenv(PathEnvVar)
	if resultPath == "" {
		// path env var not set so we'll assume we don't need to write the result to a file
		return nil
	}

	file, err := os.Create(resultPath)
	if err != nil {
		return fmt.Errorf("Error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(result)
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %w", err)
	}

	return nil
}

// Reads the result file stored by the lazygit test, and deletes the file
func ReadResult(resultPath string) (IntegrationTestResult, error) {
	file, err := os.Open(resultPath)
	if err != nil {
		return IntegrationTestResult{}, fmt.Errorf("Error reading file: %w", err)
	}

	decoder := json.NewDecoder(file)
	var result IntegrationTestResult
	err = decoder.Decode(&result)
	if err != nil {
		file.Close()
		return IntegrationTestResult{}, fmt.Errorf("Error decoding JSON: %w", err)
	}

	file.Close()
	// _ = os.Remove(resultPath)

	return result, nil
}

func SetResultPathEnvVar(cmd *exec.Cmd, resultPath string) {
	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("%s=%s", PathEnvVar, resultPath),
	)
}

func GetResultPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("lazygit_result_%s.json", generateRandomString(10)))
}

func generateRandomString(length int) string {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		panic(fmt.Sprintf("Could not generate random string: %s", err))
	}

	randomString := base64.URLEncoding.EncodeToString(buffer)
	return randomString[:length]
}
