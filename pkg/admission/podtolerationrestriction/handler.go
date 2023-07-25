/**
 * Mutate pod manifest field nodeSelector with proper key values so the pod can be scheduled to
 * designate nodes.
 */
package podtolerationrestriction

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/liangrog/admission-webhook-server/pkg/admission/admit"
	"github.com/liangrog/admission-webhook-server/pkg/utils"
	admissionV1 "k8s.io/api/admission/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	handlerName = "PodTolerationRestriction"
	// Path for kube api server to call
	ENV_POD_TOLERATION_RESTRICTION_PATH = "POD_TOLERATION_RESTRICTION_PATH"
	// Configuration for specify tolerations to namespace.
	ENV_POD_TOLERATION_RESTRICTION_CONFIG = "POD_TOLERATION_RESTRICTION_CONFIG"
)

var (
	podResource = metaV1.GroupVersionResource{Version: "v1", Resource: "pods"}
)

// Register handler to server
func Register(mux *http.ServeMux) {
	// Sub path
	path := filepath.Join(
		admit.GetBasePath(),
		utils.GetEnvVal(ENV_POD_TOLERATION_RESTRICTION_PATH,
			"pod-toleration-restriction"),
	)

	mux.Handle(path, admit.AdmitFuncHandler(handler))

	log.Printf("%s registered using path %s", handlerName, path)
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

// Get configuration map
func getConfiguredTolerationsMap() (map[string][]coreV1.Toleration, error) {
	// Don't process if no configuration is set
	s := os.Getenv(ENV_POD_TOLERATION_RESTRICTION_CONFIG)
	tolerations := map[string][]coreV1.Toleration{}
	if s != "" {
		if err := json.Unmarshal([]byte(s), &tolerations); err != nil {
			return nil, err
		}
	}
	return tolerations, nil

}
