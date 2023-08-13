package bootstrap

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeconfig "k8s.io/client-go/tools/clientcmd/api/v1"
	kubelet "k8s.io/kubelet/config/v1beta1"
)

func DefaultKubeConfig() kubeconfig.Config {
	return kubeconfig.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters: []kubeconfig.NamedCluster{kubeconfig.NamedCluster{
			Name: "default",
			Cluster: kubeconfig.Cluster{
				Server:               "",
				CertificateAuthority: ClusterCACertPath,
			},
		}},
		Contexts: []kubeconfig.NamedContext{kubeconfig.NamedContext{
			Name: "default",
			Context: kubeconfig.Context{
				Cluster:  "default",
				AuthInfo: "default",
			},
		}},
		CurrentContext: "default",
	}
}

func DefaultKubeletConfiguration() kubelet.KubeletConfiguration {
	return kubelet.KubeletConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubelet.SchemeGroupVersion.Identifier(),
			Kind:       "KubeletConfiguration",
		},
		Authentication: kubelet.KubeletAuthentication{
			X509: kubelet.KubeletX509Authentication{
				ClientCAFile: ClusterCACertPath,
			},
		},
		ShutdownGracePeriod: metav1.Duration{
			Duration: 30 * time.Second,
		},
		ShutdownGracePeriodCriticalPods: metav1.Duration{
			Duration: 10 * time.Second,
		},
		StaticPodPath: "/etc/kubernetes/manifests",
	}
}

func DefaultCredentialProviderConfig() kubelet.CredentialProviderConfig {
	return kubelet.CredentialProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubelet.SchemeGroupVersion.Identifier(),
			Kind:       "CredentialProviderConfig",
		},
	}
}
