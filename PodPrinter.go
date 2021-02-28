package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)


// Derived from kubernetes-master/pkg/printers/internalversion/printers.go
func getDescription(pod v1.Pod) string {
	var reason string
	initializing := false
	for i := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[i]
		// restarts += int(container.RestartCount)
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			// initialization is failed
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init:Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init:ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = "Init:" + container.State.Waiting.Reason
			initializing = true
		default:
			reason = fmt.Sprintf("Init:%d/%d", i, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}
	if !initializing {
		// restarts := 0
		// hasRunning := false
		for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
			container := pod.Status.ContainerStatuses[i]

			// restarts += int(container.RestartCount)
			if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
				reason = container.State.Waiting.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
				reason = container.State.Terminated.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else if container.Ready && container.State.Running != nil {
				// hasRunning = true
				// readyContainers++
			}
		}

		// change pod status back to "Running" if there is at least one container still reporting as "Running" status
		// if reason == "Completed" && hasRunning {
		// 	if hasPodReadyCondition(pod.Status.Conditions) {
		// 		reason = "Running"
		// 	} else {
		// 		reason = "NotReady"
		// 	}
		// }
	}

	if pod.DeletionTimestamp != nil {
		reason = "Terminating"
	}

	return reason
}
