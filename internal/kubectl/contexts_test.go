package kubectl

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"

	"github.com/adalbertjnr/kcmgr/internal/models"
)

var originalExecCommand = execCommand

func mockExecCommand(output string) func(string, ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		return fakeExecCommand(output)
	}
}

func fakeExecCommand(output string) *exec.Cmd {
	return exec.Command("echo", output)
}

func restoreExecCommand() {
	execCommand = originalExecCommand
}

func TestCurrentContext(t *testing.T) {
	execCommand = mockExecCommand("my-current-context")
	defer restoreExecCommand()

	result, err := CurrentContext()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != "my-current-context" {
		t.Errorf("Expected 'my-current-context', got '%s'", result)
	}
}

func TestKubernetesContexts(t *testing.T) {
	mockJSON := `
	{
		"contexts": [
			{
				"name": "ctx-1",
				"context": {
					"cluster": "cluster-1",
					"user": "user-1"
				}
			},
			{
				"name": "ctx-2",
				"context": {
					"cluster": "cluster-2",
					"user": "user-2"
				}
			}
		]
	}`

	execCommand = mockExecCommand(mockJSON)
	defer restoreExecCommand()

	contexts, err := KubernetesContexts()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(contexts) != 2 {
		t.Errorf("Expected 2 contexts, got %d", len(contexts))
	}

	ctx, ok := contexts[0].(*models.Context)
	if !ok {
		t.Errorf("Expected type *Context, got %T", contexts[0])
	}
	if ctx.Name != "ctx-1" || ctx.Context.Cluster != "cluster-1" {
		t.Errorf("Unexpected context data: %+v", ctx)
	}
}

func TestGetRawContext(t *testing.T) {
	clusterJSON := `{"name":"cluster-1","cluster":{"certificate-authority-data":"fake-cert","server":"https://127.0.0.1"}}`
	userJSON := `{"name":"user-1","user":{"token":"fake-token"}}`
	contextJSON := `{"name":"ctx-1","context":{"cluster":"cluster-1","user":"user-1"}}`

	mockOutputs := []string{clusterJSON, userJSON, contextJSON}
	currentCall := 0

	execCommand = func(command string, args ...string) *exec.Cmd {
		output := mockOutputs[currentCall]
		currentCall++
		return fakeExecCommand(output)
	}
	defer restoreExecCommand()

	result, err := GetRawContext("cluster-1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, `"name": "cluster-1"`) || !strings.Contains(result, `"name": "user-1"`) {
		t.Errorf("Expected cluster and user info in output JSON: %s", result)
	}
}

func TestContextUnmarshal(t *testing.T) {
	raw := `{"name":"ctx-1","context":{"cluster":"cluster-1","user":"user-1"}}`
	var ctx models.Context
	err := json.Unmarshal([]byte(raw), &ctx)
	if err != nil {
		t.Fatalf("Failed to unmarshal Context: %v", err)
	}

	if ctx.Name != "ctx-1" || ctx.Context.Cluster != "cluster-1" || ctx.Context.User != "user-1" {
		t.Errorf("Unexpected parsed context: %+v", ctx)
	}
}
