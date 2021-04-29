package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"
	klog "k8s.io/klog/v2"
)

const HacheckPodName = "hacheck"

// fetch running pods from kubelet and update pods registration list
func fetchPods(ctx context.Context, out chan []ServiceRegistration) {
	var oldRegistrations, newRegistrations []ServiceRegistration
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			klog.Info("fetchPods stopped")
			return
		case <-ticker.C:
			startTime := time.Now().UnixNano()
			resp, err := http.Get("http://127.0.0.1:10255/pods")
			if err != nil {
				klog.Errorf("fetching pods failed: %v", err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				klog.Errorf("fetching pods bad response: %v", resp.StatusCode)
				continue
			}

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				klog.Errorf("reading body failed: %v", err)
				continue
			}
			loadedTime := time.Now().UnixNano()
			klog.Infof(
				"read %v bytes in %vs",
				len(bodyBytes),
				float64(loadedTime-startTime)/float64(time.Second),
			)

			var podList corev1.PodList
			err = json.Unmarshal(bodyBytes, &podList)
			if err != nil {
				klog.Errorf("unmarshaling body failed: %v", err)
				continue
			}
			parsedTime := time.Now().UnixNano()

			klog.Infof(
				"loaded %v pods in %vs",
				len(podList.Items),
				float64(parsedTime-loadedTime)/float64(time.Second),
			)

			newRegistrations = []ServiceRegistration{}
			for _, pod := range podList.Items {
				if pod.Status.Phase != corev1.PodRunning {
					continue
				}
				podRegsJson, ok := pod.Annotations["smartstack_registrations"]
				if !ok {
					continue
				}

				var podRegs []string
				err := json.Unmarshal([]byte(podRegsJson), &podRegs)
				if err != nil {
					klog.Errorf(
						"pod %v/%v smartstack_registrations failed to load: %v, raw json: %+v",
						pod.Namespace,
						pod.Name,
						err,
						podRegsJson,
					)
					continue
				}

				var port int32
				for _, cont := range pod.Spec.Containers {
					// TODO: use instance name?
					if cont.Name != HacheckPodName {
						port = cont.Ports[0].ContainerPort
						break
					}
				}
				service := pod.Labels["paasta.yelp.com/service"]
				instance := pod.Labels["paasta.yelp.com/instance"]
				podIP := pod.Status.PodIP

				for _, reg := range podRegs {
					newRegistrations = append(newRegistrations, ServiceRegistration{
						Service:      service,
						Instance:     instance,
						PodNode:      pod.Spec.NodeName,
						PodNs:        pod.Namespace,
						PodName:      pod.Name,
						PodIP:        podIP,
						Port:         port,
						Registration: reg,
					})
				}
			}

			if registrationsEqual(oldRegistrations, newRegistrations) {
				klog.V(10).Info("pods registrations did not change")
				continue
			}

			klog.Infof("pods registrations updated: %+v", newRegistrations)
			oldRegistrations = newRegistrations
			out <- newRegistrations
		}
	}
}
