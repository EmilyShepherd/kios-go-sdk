package bootstrap

import (
	"fmt"
	"os"

	"k8s.io/klog/v2"
)

const ClusterCADir = "/etc/kubernetes/pki"
const ClusterCACertPath = ClusterCADir + "/ca.crt"
const ClusterCAKeyPath = ClusterCADir + "/ca.key"

// Represents a certificate that can have both public and private parts
type Cert struct {
	Cert []byte
	Key  []byte
}

// Saves a byte stream to the given disk if it is not empty. This is
// used by SaveClusterCA as the private key is optional.
func saveIfPresent(data []byte, path string, mode os.FileMode) error {
	if len(data) > 0 {
		if err := os.WriteFile(path, data, mode); err != nil {
			return fmt.Errorf("Could not write cluster CA file to disk: %s", err)
		}

		klog.Infof("Cluster CA file written to disk: %s", path)
	}

	return nil
}

// Saves the cluster's CA data to the correct location on disk
func (b *Bootstrap) SaveClusterCA() error {
	ca := b.Provider.GetClusterCA()

	if err := os.MkdirAll(ClusterCADir, 0755); err != nil {
		return fmt.Errorf("Could not create Cluster CA Directory: %s", err)
	}

	if err := saveIfPresent(ca.Cert, ClusterCACertPath, 0644); err != nil {
		return err
	}
	if err := saveIfPresent(ca.Key, ClusterCAKeyPath, 0600); err != nil {
		return err
	}

	return nil
}
