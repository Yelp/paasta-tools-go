package acceptance

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/subosito/gotenv"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func buildEnv(makefile string) error {
	cmd := exec.Command("make", "local-env", "-f", makefile)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.Printf("%q", out.String())
	log.Printf("%q", stderr.String())
	if err != nil {
		log.Fatal(err)
		return err
	}
	env, err := gotenv.StrictParse(&out)
	if err != nil {
		return err
	}

	for key, val := range env {
		fmt.Println(fmt.Sprintf("%s: %s", key, val))
		if _, present := os.LookupEnv(key); !present {
			os.Setenv(key, val)
		}
	}
	return nil
}

// KubernetesCluster creates a kubernetes cluster using a Makefile
func KubernetesCluster(makefile string) (*kubernetes.Clientset, error) {
	err := buildEnv(makefile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	cmd := exec.Command("make", "local-cluster", "-f", makefile)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	log.Printf("%q", out.String())
	log.Printf("%q", stderr.String())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
