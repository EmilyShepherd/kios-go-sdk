package bootstrap

import (
	"fmt"
	"os"
)

// This file is read during kiOS' boot and is used to set node labels
const NodeLabelsPath = "/etc/kubernetes/node-labels"

// Generates the node-labels file, which is used by kiOS to set the
// labels which kubelet should register itself with
func (b *Bootstrap) SaveNodeLabels() error {
	f, err := os.OpenFile(NodeLabelsPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	for key, value := range b.Provider.GetNodeLabels() {
		f.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	return f.Close()
}
