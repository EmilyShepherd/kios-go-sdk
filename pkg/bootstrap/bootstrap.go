package bootstrap

import (
	"fmt"
	"k8s.io/klog/v2"

	"github.com/EmilyShepherd/kios-go-sdk/pkg/socket"
)

type Bootstrap struct {
	// An instance of a BootstrapInformationProvider is required. This
	// will be called to provide the information used to bootstrap the
	// node.
	Provider BootstrapInformationProvider

	// A list of binaries to copy across.
	// These are always processed before anything else, to allow time for
	// background copying on slow harddrives.
	Binaries []string
}

// Run a standard node bootstrap. Alternatively you may call each step
// individually if you have another order / way that you would like to
// bootstrap the node.
func (b *Bootstrap) Run() {
	c, _ := b.CopyBinaries()

	b.Provider.Init()

	systemSocket, err := socket.NewSystemSocket()
	if err != nil {
		fmt.Printf("Could not open a connection to the system socket: %s\n", err)
	}

	b.SaveClusterCA()
	b.SaveCredentialProviderConfig()
	b.SaveHostnameFile()
	b.SaveNodeLabels()
	b.SaveKubeConfig()
	b.SaveKubeletConfiguration()

	crioUpdated, err := b.SaveCrioConfiguration()
	if err != nil {
		klog.Errorf("Could not save crio configuration: %s", err)
	} else if crioUpdated {
		klog.Warning("MetadataConfiguration reconfigures the container runtime. Crio will be restarted")
		systemSocket.SendCmd(socket.CmdRestartCrio)
	}

	b.WaitForBinaries(c)

	systemSocket.SendCmd(socket.CmdRestartKubelet)
}
