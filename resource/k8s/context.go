package k8s

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// ContextRes represents a kubernetes clusters configurations.
type ContextRes interface {
	Res
	Switch(n string) error
}

// NamedContext represents a named cluster context.
type NamedContext struct {
	Name    string
	Context *api.Context
}

// CurrentCluster return active cluster name
func (c *NamedContext) CurrentCluster() string {
	return conn.getClusterName()
}

// Context represents a Kubernetes Context.
type Context struct{}

// NewContext returns a new Context.
func NewContext() Res {
	return &Context{}
}

// Get a Context.
func (*Context) Get(_, n string) (interface{}, error) {
	return &NamedContext{Name: n, Context: conn.apiConfigOrDie().Contexts[n]}, nil
}

// List all Contexts in a given namespace
func (*Context) List(string) (Collection, error) {
	conn := conn.apiConfigOrDie()

	cc := make([]interface{}, 0, len(conn.Contexts))
	for k, v := range conn.Contexts {
		cc = append(cc, &NamedContext{k, v})
	}
	return cc, nil
}

// Delete a Context
func (*Context) Delete(_, n string) error {
	conn := conn.apiConfigOrDie()

	if conn.CurrentContext == n {
		return fmt.Errorf("trying to delete your current context %s", n)
	}

	acc := clientcmd.NewDefaultPathOptions()

	delete(conn.Contexts, n)
	return clientcmd.ModifyConfig(acc, *conn, true)
}

// Switch cluster Context.
func (*Context) Switch(n string) error {
	conn := conn.apiConfigOrDie()

	conn.CurrentContext = n
	acc := clientcmd.NewDefaultPathOptions()
	return clientcmd.ModifyConfig(acc, *conn, true)
}
