package bootstrap

import (
	"fmt"
	"os"
)

const HostnameFilePath = "/etc/hostname"

// Assuming this is running in the kubelet's bootstrap run, it is
// acceptable to simply write the desired hostname to the host's
// /etc/hostname file. Init will pick this ip and auto set the hostname
// before restarting the kubelet.
func (b *Bootstrap) SaveHostnameFile() error {
	if err := os.WriteFile(HostnameFilePath, []byte(b.Provider.GetHostname()), 0644); err != nil {
		return fmt.Errorf("Could not write hostname file: %s", err)
	}

	return nil
}
