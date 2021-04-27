package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"k8s.io/klog/v2"
)

// check host for services configured via other means and update registration list
func fetchLocal(ctx context.Context, out chan []ServiceRegistration, servicesDir string) {
	var oldRegistrations, newRegistrations []ServiceRegistration
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			klog.Info("fetchLocal stopped")
			return
		case <-ticker.C:
			newRegistrations = []ServiceRegistration{}
			err := filepath.Walk(servicesDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return fmt.Errorf("error accessing %v: %v", path, err)
				}
				if info.IsDir() {
					return nil
				}

				fh, err := os.Open(path)
				if err != nil {
					// ignore broken symlinks
					return nil
				}

				serviceInfo := struct {
					Namespaces []string `json:"namespaces"`
				}{}
				bytes, err := ioutil.ReadAll(fh)
				if err != nil {
					return fmt.Errorf("error reading %v: %v", path, err)
				}
				err = json.Unmarshal(bytes, &serviceInfo)
				if err != nil {
					return fmt.Errorf("error parsing %v: %v", path, err)
				}
				if serviceInfo.Namespaces == nil {
					klog.Infof("service has no namespaces: %v", path)
					return nil
				}
				klog.Infof("local service: %v %+v\n", info.Name(), serviceInfo.Namespaces)
				return nil
			})
			if err != nil {
				klog.Error(err)
				continue
			}
			if registrationsEqual(oldRegistrations, newRegistrations) {
				continue
			}
			out <- newRegistrations
		}
	}
}
