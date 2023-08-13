package bootstrap

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slices"
	kubeconfig "k8s.io/client-go/tools/clientcmd/api/v1"
	"k8s.io/klog/v2"
	kubelet "k8s.io/kubelet/config/v1beta1"
	"sigs.k8s.io/yaml"
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

	kubeConfigStr, err := yaml.Marshal(kubeConfig)

	if err != nil {
		return fmt.Errorf("Could not marshal KubeConfig YAML: %s", err)
	}

	if err = os.WriteFile(KubeletKubeconfigPath, kubeConfigStr, 0600); err != nil {
		return fmt.Errorf("Could not write Kubeconfig to disk: %s", err)
	}

	klog.Infof("KubeConfig written to disk: %s", KubeletKubeconfigPath)

	return nil
}

// Reads the given file from disk and unmarshals it as YAML
func yamlFromFile(filename string, obj interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Could not open file %s: %s", filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Could not read file %s: %s", filename, err)
	}

	if err := yaml.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("Could not parse YAML from file %s: %s", filename, err)
	}

	return nil
}

// Loads the template kubeconfig file from disk, adds the relavent
// settings to it, before remarshalling it as YAML and saving it back to
// disk
func (b *Bootstrap) SaveKubeletConfiguration() error {
	kubeletConfig := DefaultKubeletConfiguration()

	if err := yamlFromFile(KubeletConfigurationPath, &kubeletConfig); err != nil {
		klog.Warning(err.Error())
	}

	kubeletConfig = b.Provider.GetKubeletConfiguration(kubeletConfig)

	kubelet, _ := yaml.Marshal(&kubeletConfig)
	os.WriteFile(KubeletConfigurationPath, kubelet, 0644)

	klog.Infof("Kubelet Configuration written to disk: %s", KubeletConfigurationPath)

	return nil
}

// Creates the credential provider configuration file for image
// credentials
func (b *Bootstrap) SaveCredentialProviderConfig() error {
	config := DefaultCredentialProviderConfig()

	if err := yamlFromFile(CredentialProviderConfigPath, &config); err != nil {
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

	providerConfig, _ := yaml.Marshal(&config)
	os.WriteFile(CredentialProviderConfigPath, providerConfig, 0600)

	klog.Infof("Credential Provider Config written to disk: %s", CredentialProviderConfigPath)

	return nil
}
