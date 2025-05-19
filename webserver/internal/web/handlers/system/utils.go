package system

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"

	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
)

var (
	Scheme = runtime.NewScheme()
	_      = clusterctlv1.AddToScheme(Scheme) // Register Cluster API types
	_      = clusterv1.AddToScheme(Scheme)    // Register Cluster API types
	_      = capv.AddToScheme(Scheme)         // Register Cluster API types
)

// HandleError write down an error with code to the writer response.
func HandleError(w http.ResponseWriter, code int, err error) (hasError bool) {
	hasError = err != nil
	if hasError {
		writeError(w, code, err)
	}
	return hasError
}

// writeError write down an error with code to the writer response.
func writeError(w http.ResponseWriter, code int, err error) {
	log.Println("ERROR: ", err)
	http.Error(w, err.Error(), code)
}

// WriteResponse write the response byte input on writer.
func WriteResponse(w http.ResponseWriter, object any) error {
	response, err := convertObject(object)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(string(response))
	if _, err := w.Write(response); err != nil {
		return err
	}
	return nil
}

// convertObject marshal to a generic object on a []byte return.
func convertObject(object any) (response []byte, err error) {
	if response, err = json.Marshal(&object); err != nil {
		fmt.Println(err, "ERROR")
		return make([]byte, 0), err
	}
	return response, nil
}
