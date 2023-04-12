package kubernetes

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

type Client struct {
	Client     kubernetes.Interface
	RestClient rest.Interface
	Config     *rest.Config
}

func (c *Client) GetConfig() *rest.Config {
	return c.Config
}

func (c *Client) GetClient() kubernetes.Interface {
	return c.Client
}

func (c *Client) GetRestClient() rest.Interface {
	return c.RestClient
}

func NewClient(kubecontext string, kubeconfig string) (*Client, error) {

	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: kubecontext,
		})
	// create the clientset
	c, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	c.APIPath = "/api"
	c.GroupVersion = &scheme.Scheme.PrioritizedVersionsForGroup("")[0]
	c.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}

	restClient, err := rest.RESTClientFor(c)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:     clientSet,
		RestClient: restClient,
		Config:     c,
	}, nil
}
