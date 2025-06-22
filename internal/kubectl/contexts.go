package kubectl

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/adalbertjnr/kcmgr/internal/models"
	"github.com/charmbracelet/bubbles/list"
)

var execCommand = exec.Command

type Cluster struct {
	Name    string `json:"name"`
	Cluster struct {
		Certificate string `json:"certificate-authority-data"`
		Server      string `json:"server"`
	} `json:"cluster"`
}

type User struct {
	Name string      `json:"name"`
	User interface{} `json:"user"`
}

func KubernetesContexts() ([]list.Item, error) {
	contextsCommand := execCommand("kubectl", "config", "view", "-o", "json")
	contextsOutput, err := contextsCommand.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var raw struct {
		Contexts []models.Context `json:"contexts"`
	}
	if err := json.Unmarshal(contextsOutput, &raw); err != nil {
		return nil, err
	}

	items := make([]list.Item, len(raw.Contexts))
	for i := range raw.Contexts {
		items[i] = &raw.Contexts[i]
	}

	return items, nil
}

type RawContext struct {
	Cluster Cluster
	Context models.Context
	User    User
}

func GetRawContext(clusterName string) (string, error) {
	clusterQuery := fmt.Sprintf(`jsonpath={.clusters[?(@.name=="%s")]}`, clusterName)
	contextQuery := fmt.Sprintf(`jsonpath={.contexts[?(@.context.cluster=='%s')]}`, clusterName)
	userQuery := fmt.Sprintf(`jsonpath={.users[?(@.name=='%s')]}`, clusterName)

	var cluster Cluster
	if err := runKubectlJSONPath(clusterQuery, &cluster); err != nil {
		return "", err
	}

	var user User
	if err := runKubectlJSONPath(userQuery, &user); err != nil {
		return "", err
	}

	var context models.Context
	if err := runKubectlJSONPath(contextQuery, &context); err != nil {
		return "", err
	}

	rc := RawContext{
		Cluster: cluster,
		Context: context,
		User:    user,
	}

	rawContextString, err := json.MarshalIndent(&rc, "", "  ")
	if err != nil {
		return "", err
	}

	return string(rawContextString), nil
}

func runKubectlJSONPath(query string, out any) error {
	queryCommand := execCommand("kubectl", "config", "view", "--raw", "-o", query)
	output, err := queryCommand.CombinedOutput()
	if err != nil {
		return err
	}
	return json.Unmarshal(output, out)
}

func CurrentContext() (string, error) {
	contextCommand := execCommand("kubectl", "config", "current-context")
	context, err := contextCommand.CombinedOutput()
	if err != nil {
		return "", err
	}

	trimmedSpaceContext := strings.TrimSpace(string(context))
	return trimmedSpaceContext, nil
}

func SetKubernetesContext(context string) error {
	useContextCommand := execCommand("kubectl", "config", "use-context", context)
	return useContextCommand.Run()
}

func DeleteKubernetesContext(context string) error {
	deleteContextCommand := execCommand("kubectl", "config", "delete-context", context)
	return deleteContextCommand.Run()
}
