package main

import (
	"context"
	"os"
	"os/signal"

	"k8s.io/klog/v2"
)

func signalLoop(ctx context.Context, cancel context.CancelFunc) {
	stopping := false
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	for s := range signalCh {
		if stopping {
			klog.Warningf("caught %v, still stopping...", s)
			continue
		}

		klog.Infof("caught %v, stopping...", s)
		cancel()
		stopping = true
	}
}
