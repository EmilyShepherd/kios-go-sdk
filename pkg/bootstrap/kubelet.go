package bootstrap

import (
	"github.com/EmilyShepherd/kios-go-sdk/pkg/yaml"
	"golang.org/x/exp/slices"
	kubeconfig "k8s.io/client-go/tools/clientcmd/api/v1"
	"k8s.io/klog/v2"
	kubelet "k8s.io/kubelet/config/v1beta1"
)

// The path where the kubelet expects its various configurs to exist.
// These are hard coded into kios' init, so cannot be changed here.
const KubeletKubeconfigPath = "/etc/kubernetes/kubelet.conf"
const KubeletConfigurationPath = "/etc/kubernetes/config.yaml"
const CredentialProviderConfigPath = "/etc/kubernetes/credential-providers.yaml"

// Generates a KubeConfig file for Kubelet, marshals it to YAML, and
// saves it
func (b *Bootstrap) SaveKubeConfig() error {
	kubeConfig := DefaultKubeConfig()
	kubeConfig.Clusters[0].Cluster.Server = b.Provider.GetClusterEndpoint()
	kubeConfig.AuthInfos = []kubeconfig.NamedAuthInfo{kubeconfig.NamedAuthInfo{
		Name:     "default",
		AuthInfo: b.Provider.GetClusterAuthInfo(),
	}}

	return yaml.YamlToFile(kubeConfig, KubeletKubeconfigPath, 0600)
}

// Loads the template kubeconfig file from disk, adds the relavent
// settings to it, before remarshalling it as YAML and saving it back to
// disk
func (b *Bootstrap) SaveKubeletConfiguration() error {
	kubeletConfig := DefaultKubeletConfiguration()

	if err := yaml.YamlFromFile(KubeletConfigurationPath, &kubeletConfig); err != nil {
		klog.Warning(err.Error())
	}

	kubeletConfig = b.Provider.GetKubeletConfiguration(kubeletConfig)

	return yaml.YamlToFile(kubeletConfig, KubeletConfigurationPath, 0644)
}

// Creates the credential provider configuration file for image
// credentials
func (b *Bootstrap) SaveCredentialProviderConfig() error {
	config := DefaultCredentialProviderConfig()

	if err := yaml.YamlFromFile(CredentialProviderConfigPath, &config); err != nil {
		klog.Warning(err.Error())
	}

	for _, provider := range b.Provider.GetCredentialProviders() {
		idx := slices.IndexFunc(
			config.Providers,
			func(p kubelet.CredentialProvider) bool {
				return p.Name == provider.Name
			},
		)

		if idx == -1 {
			config.Providers = append(config.Providers, provider)
		} else {
			config.Providers[idx] = provider
		}
	}

	return yaml.YamlToFile(config, CredentialProviderConfigPath, 0600)
}
