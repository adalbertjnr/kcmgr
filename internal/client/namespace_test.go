package client

import (
	"context"
	"testing"
	"time"

	"github.com/adalbertjnr/kcmgr/internal/models"
	"github.com/adalbertjnr/kcmgr/internal/util"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type fakeNamespaceLister struct {
	outputList *corev1.NamespaceList
	outputErr  error
}

func (f *fakeNamespaceLister) List(ctx context.Context, opts v1.ListOptions) (*corev1.NamespaceList, error) {
	return f.outputList, f.outputErr
}

func TestGetNamespaceByContext(t *testing.T) {

	now := time.Now().UTC()
	defaultNamespaceTime := v1.NewTime(now)
	kubeSystemTime := v1.NewTime(now.Add(time.Hour))
	argocdTime := v1.NewTime(now.Add(-5 * time.Hour))

	expectedDefaultTime, err := util.ParseTime(defaultNamespaceTime.String())
	assert.NoError(t, err)
	expectedKubeSystemTime, err := util.ParseTime(kubeSystemTime.String())
	assert.NoError(t, err)
	expectedArgoCDTime, err := util.ParseTime(argocdTime.String())
	assert.NoError(t, err)

	expected := []models.Namespace{
		{Name: "default", Age: expectedDefaultTime},
		{Name: "kube-system", Age: expectedKubeSystemTime},
		{Name: "argocd", Age: expectedArgoCDTime},
	}

	fake := &fakeNamespaceLister{
		outputList: &corev1.NamespaceList{
			Items: []corev1.Namespace{
				{ObjectMeta: v1.ObjectMeta{Name: expected[0].Name, CreationTimestamp: defaultNamespaceTime}},
				{ObjectMeta: v1.ObjectMeta{Name: expected[1].Name, CreationTimestamp: kubeSystemTime}},
				{ObjectMeta: v1.ObjectMeta{Name: expected[2].Name, CreationTimestamp: argocdTime}},
			},
		},
	}

	namespaces, err := toModelNamespace(fake.outputList)

	assert.NoError(t, err)
	assert.Equal(t, expected[0].Name, namespaces[0].Name)
	assert.Equal(t, expected[1].Name, namespaces[1].Name)
	assert.Equal(t, expected[2].Name, namespaces[2].Name)
	assert.Equal(t, expectedDefaultTime, namespaces[0].Age)
	assert.Equal(t, expectedKubeSystemTime, namespaces[1].Age)
	assert.Equal(t, expectedArgoCDTime, namespaces[2].Age)
}
