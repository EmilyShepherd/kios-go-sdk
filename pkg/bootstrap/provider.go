package bootstrap

import (
	kubeconfig "k8s.io/client-go/tools/clientcmd/api/v1"
	kubelet "k8s.io/kubelet/config/v1beta1"
)

// The BoostrapInformationProvider interface is expected to be
// implemented by downstream projects.
//
// Eg, kios-aws might use this interface to provide required config
// values by pulling information from the EC2 metadata service. Another
// provider may choose to load information from UEFI variables etc...
type BootstrapInformationProvider interface {
	Init() error
	GetClusterCA() Cert
	GetCredentialProviders() []kubelet.CredentialProvider
	GetHostname() string
	GetNodeLabels() map[string]string
	GetClusterEndpoint() string
	GetClusterAuthInfo() kubeconfig.AuthInfo
	GetContainerRuntimeConfiguration() ContainerRuntimeConfiguration
	GetKubeletConfiguration(kubelet.KubeletConfiguration) kubelet.KubeletConfiguration
}
