// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authorizaton

import (
	"fmt"
	"strings"
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

// func cleanupAuthorEA() {
// 	util.Log.Info("Cleanup")
// 	sleep := examples.Sleep{"foo"}
// 	httpbin := examples.Httpbin{"foo"}
// 	sleep.Uninstall()
// 	httpbin.Uninstall()
// 	time.Sleep(time.Duration(20) * time.Second)
// }

func TestAuthorEA(t *testing.T) {
	// defer cleanupAuthorEA()
	defer util.RecoverPanic(t)

	// util.Log.Info("Authorization with JWT Token")
	// sleep := examples.Sleep{"foo"}
	// sleep.Install()
	// httpbin := examples.Httpbin{"foo"}
	// httpbin.Install()

	// sleepPod, err := util.GetPodName("foo", "app=sleep")
	// util.Inspect(err, "Failed to get sleep pod name", "", t)
	// cmd := fmt.Sprintf(`curl http://httpbin.foo:8000/ip -sS -o /dev/null -w "%%{http_code}\n"`)
	// msg, err := util.PodExec("foo", sleepPod, "sleep", cmd, true)
	// util.Inspect(err, "Failed to get response", "", t)
	// if !strings.Contains(msg, "200") {
	// 	util.Log.Errorf("Verify setup -- Unexpected response code: %s", msg)
	// } else {
	// 	util.Log.Infof("Success. Get expected response: %s", msg)
	// }

	// util.Log.Info("Deploy the External Authorizer")
	// util.KubeApplyContents("foo", ExternalAuthorizer)
	// time.Sleep(time.Duration(50) * time.Second)

	// util.Log.Info("Verfiy the sample external authorizer is up and running")

	// extAuthPod, err := util.GetPodName("foo", "app=ext-authz")
	// util.Inspect(err, "Failed to get ext-authz pod name", "", t)
	// // util.Shell("kubectl logs %s -n foo -c ext-authz", extAuthPod)

	// extAuthPodverf, err := util.Shell("oc logs %s -n foo -c ext-authz", extAuthPod)
	// if err != nil {
	// 	t.Fatalf("external autorization is not running: %v", err)
	// }
	// extAuthPodverf = strings.Trim(extAuthPodverf, "'")

	// ext, err := util.Shell("oc get smcp -n %s -o jsonpath={.items..metadata.name}", meshNamespace)
	// util.Inspect(err, "Failed to smcp name", "", t)
	// util.Shell("oc get configmap istio-%s -n %s", ext, meshNamespace)

	// // oc get configmap istio-basic -o jsonpath={.data} -n istio-system  --> configmap jsonpath
	// util.Log.Info("Edit configmap for two external providers")
	// util.Shell(`kubectl patch -n %s configmap/istio-basic --type merge -p '{"data": {"mesh": "# Add the following content to define the external authorizers.\nextensionProviders:\n- name: \"sample-ext-authz-grpc\"\n  envoyExtAuthzGrpc:\n    service: \"ext-authz.foo.svc.cluster.local\"\n    port: \"9000\"\n- name: \"sample-ext-authz-http\"\n  envoyExtAuthzHttp:\n    service: \"ext-authz.foo.svc.cluster.local\"\n    port: \"8000\"\n    includeRequestHeadersInCheck: [\"x-ext-authz\"]"}}'`, meshNamespace)

	// util.Log.Info("Enable with external authorization")
	// util.KubeApplyContents("foo", ExternalRoute)

	// util.Log.Info("Verfiy the header with deny server")

	sleepPod, err := util.GetPodName("foo", "app=sleep")
	// util.Inspect(err, "Failed to get sleep pod name", "", t)
	// deny := fmt.Sprintf(`curl "http://httpbin.foo:8000/headers" -H "x-ext-authz: deny" -s`)
	// util.Shell("kubectl exec %s -c sleep -n foo -- %s", sleepPod, deny)

	// util.Log.Info("Verfiy the header with allow server")

	// util.Inspect(err, "Failed to get sleep pod name", "", t)
	// allow := fmt.Sprintf(`curl "http://httpbin.foo:8000/headers" -H "x-ext-authz: allow" -s`)
	// util.Shell("kubectl exec %s -c sleep -n foo -- %s", sleepPod, allow)

	util.Log.Info("Verify request with a JWT includes group1 claim")
	cmd := fmt.Sprintf(`curl "http://httpbin.foo:8000/ip" -s -o /dev/null -w "%{http_code}\n"`, "foo")
	msg, err := util.PodExec("foo", sleepPod, "sleep", cmd, true)
	util.Inspect(err, "Failed to get response", "", t)
	if !strings.Contains(msg, "200") {
		util.Log.Errorf("Verify request with JWT group1 claim Unexpected response: %s", msg)
		t.Errorf("Verify request with JWT group1 claim Unexpected response: %s", msg)
	} else {
		util.Log.Infof("Success. Get expected response: %s", msg)
	}

}
