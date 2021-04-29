package main

import (
	"context"
	"time"

	"k8s.io/klog/v2"
)

// merge pods, local and aws, calculate diff from previous frame
// and write differences
func writeAWS(ctx context.Context, podsch, localch, awsch chan []ServiceRegistration) {
	var pods, local, aws []ServiceRegistration
	var podsDirty, localDirty, awsDirty bool
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		klog.Info("selecting")
		select {
		case <-ctx.Done():
			klog.Info("writeAWS stopped")
			return
		case pods = <-podsch:
			klog.Info("pods updated")
			podsDirty = true
		case local = <-localch:
			klog.Info("local updated")
			localDirty = true
		case aws = <-awsch:
			klog.Info("aws updated")
			awsDirty = true
		case <-ticker.C:
			if podsDirty {
				klog.Infof("new pods=%+v\n", pods)
			}
			if localDirty {
				klog.Infof("new local=%+v\n", local)
			}
			if awsDirty {
				klog.Infof("new aws=%+v\n", aws)
			}
			if !podsDirty && !localDirty && !awsDirty {
				klog.Info("nothing updated")
				continue
			}

			podsDirty = false
			localDirty = false
			awsDirty = false
		}
	}
}

// find registrations in CloudMap that match current node name
func fetchAWS(ctx context.Context, aws chan []ServiceRegistration) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			klog.Info("fetchAWS stopped")
			return
		case <-ticker.C:
			continue
		}
	}
}
