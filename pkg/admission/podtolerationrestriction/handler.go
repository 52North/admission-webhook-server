/**
 * Mutate pod manifest field nodeSelector with proper key values so the pod can be scheduled to
 * designate nodes.
 */
package podtolerationrestriction

import (
	"fmt"
	"log"
	"strings"

	"github.com/52north/admission-webhook-server/pkg/admission/admit"
	admissionV1 "k8s.io/api/admission/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	handlerName = "PodTolerationRestriction"
	// Configuration for specify tolerations to namespace.
	ENV_POD_TOLERATION_RESTRICTION_CONFIG = "POD_TOLERATION_RESTRICTION_CONFIG"
)

var (
	podResource = metaV1.GroupVersionResource{Version: "v1", Resource: "pods"}
)

func Register(ctrl admit.AdmissionController) {
	ctrl.Register(handlerName, handler)
}

// Handling pod toleration restriction request
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
	tolerationsMap, err := getConfiguredTolerationsMap()
	if err != nil {
		log.Fatal(err)
	}

	if tolerationsMap != nil {
		if tolerations, ok := tolerationsMap[req.Namespace]; ok {
			if pod.Spec.Tolerations != nil {
				for _, toleration := range tolerations {
					patches = append(patches, admit.PatchOperation{
						Op:    "add",
						Path:  "/spec/tolerations/-",
						Value: toleration,
					})
				}
			} else {
				patches = append(patches, admit.PatchOperation{
					Op:    "add",
					Path:  "/spec/tolerations",
					Value: tolerations,
				})
			}
			log.Printf("%s processed pod %s with tolerations: %s",
				handlerName,
				podName,
				fmt.Sprintf("%v", tolerations),
			)
		}
	}

	return patches, nil
}
