package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.yun.pingan.com/eks/admission-webhook/pkg/server/service"
	"github.com/gosoon/glog"
	"k8s.io/api/admission/v1beta1"
)

// admitFunc is the type we use for all of our validators and mutators
type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

// serve handles the http portion of a request prior to handing to an admit
// function
func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	glog.Info(fmt.Sprintf("handling request: %s", body))

	// The AdmissionReview that was sent to the webhook
	requestedAdmissionReview := v1beta1.AdmissionReview{}

	// The AdmissionReview that will be returned
	responseAdmissionReview := v1beta1.AdmissionReview{}

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		glog.Error(err)
		responseAdmissionReview.Response = service.ToAdmissionResponse(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = admit(requestedAdmissionReview)
	}

	// Return the same UID
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	glog.Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		glog.Error(err)
	}
	if _, err := w.Write(respBytes); err != nil {
		glog.Error(err)
	}
}

func serveClusterOperator(w http.ResponseWriter, r *http.Request) {
	serve(w, r, service.AdmitClusterOperator)
}

func servePods(w http.ResponseWriter, r *http.Request) {
	serve(w, r, service.AdmitPods)
}

func serveMutatePods(w http.ResponseWriter, r *http.Request) {
	serve(w, r, service.MutatePods)
}

func RunServer() error {
	http.HandleFunc("/eks/operator/cluster", serveClusterOperator)
	http.HandleFunc("/pods", servePods)
	http.HandleFunc("/mutating-pods", serveMutatePods)

	server := &http.Server{
		Addr: ":8080",
		//TLSConfig: configs.ConfigTLS(config),
	}
	err := server.ListenAndServe()
	return err
}
