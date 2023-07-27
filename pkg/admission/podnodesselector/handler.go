/**
 * Mutate pod manifest field nodeSelector with proper key values so the pod can be scheduled to
 * designate nodes.
 */
package podnodesselector

import (
	"fmt"
	"log"
	"strings"

	"github.com/52north/admission-webhook-server/pkg/admission/admit"
	admissionV1 "k8s.io/api/admission/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	handlerName = "PodNodesSelector"

	// Configuration for specify nodes to namespace.
	// The string format for each namespace follows https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/labels/labels.go
	// Examples:
	//   namespace:label-name=label-value,label-name=label-value;namespace:label-name=label-value
	ENV_POD_NODES_SELECTOR_CONFIG = "POD_NODES_SELECTOR_CONFIG"

	namespaceSeparator      = ";"
	namespaceLabelSeparator = ":"
)

var (
	podResource = metaV1.GroupVersionResource{Version: "v1", Resource: "pods"}
)

func Register(ctrl admit.AdmissionController) {
	ctrl.Register(handlerName, handler)
}

// Handling pod node selector request
func handler(req *admissionV1.AdmissionRequest) ([]admit.PatchOperation, error) {
	if req.Resource != podResource {
		log.Printf("Ignore admission request %s as it's not a pod resource", string(req.UID))
		return nil, nil
	}

	// Parse the Pod object.
	raw := req.Object.Raw
	pod := coreV1.Pod{}
	if _, _, err := admit.UniversalDeserializer.Decode(raw, nil, &pod); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	// Get the pod name for info
	podName := strings.TrimSpace(pod.Name + " " + pod.GenerateName)

	var patches []admit.PatchOperation

	// Get configuration
	selectors, err := getConfiguredSelectorMap()
	if err != nil {
		log.Fatal(err)
	}

	if selectors != nil {
		if labelSet, ok := selectors[req.Namespace]; ok {
			op := "replace"
			if pod.Spec.NodeSelector == nil {
				op = "add"
			}

			if labels.Conflicts(labelSet, labels.Set(pod.Spec.NodeSelector)) {
				return patches, fmt.Errorf("pod node label selector conflicts with its namespace node label selector for pod %s", podName)
			}

			podNodeSelectorLabels := labels.Merge(labelSet, labels.Set(pod.Spec.NodeSelector))

			patches = append(patches, admit.PatchOperation{
				Op:    op,
				Path:  "/spec/nodeSelector",
				Value: podNodeSelectorLabels,
			})

			log.Printf("%s processed pod %s with selectors: %s",
				handlerName,
				podName,
				fmt.Sprintf("%v", podNodeSelectorLabels),
			)
		}
	}

	return patches, nil
}
