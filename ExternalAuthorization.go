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
	"time"

	"github.com/maistra/maistra-test-tool/pkg/examples"
	"github.com/maistra/maistra-test-tool/pkg/util"
)

func cleanupAuthorEA() {
	util.Log.Info("Cleanup")
	sleep := examples.Sleep{"foo"}
	httpbin := examples.Httpbin{"foo"}
	sleep.Uninstall()
	httpbin.Uninstall()
	util.KubeDeleteContents("foo", ExternalAuthzService)
	util.KubeDeleteContents("foo", ExternalAuthzDeployment)
	util.KubeDeleteContents("foo", ExternalRoute)
	time.Sleep(time.Duration(20) * time.Second)
}

func TestAuthorEA(t *testing.T) {
	// defer cleanupAuthorEA()
	defer util.RecoverPanic(t)

	sleepPod, err := util.GetPodName("foo", "app=sleep")

	util.Log.Info("Authorization with JWT Token")
	sleep := examples.Sleep{"foo"}
	sleep.Install()
	httpbin := examples.Httpbin{"foo"}
	httpbin.Install()

	util.Inspect(err, "Failed to get sleep pod name", "", t)
	cmd := fmt.Sprintf(`curl http://httpbin.foo:8000/ip -sS -o /dev/null -w "%%{http_code}\n"`)
	msg, err := util.PodExec("foo", sleepPod, "sleep", cmd, true)
	util.Inspect(err, "Failed to get response", "", t)
	if !strings.Contains(msg, "200") {
		util.Log.Errorf("Verify setup -- Unexpected response code: %s", msg)
	} else {
		util.Log.Infof("Success. Get expected response: %s", msg)
	}

	// 1. Deploy the External Authorizer Service
	util.Log.Info("Deploy the External Authorizer")
	util.KubeApplyContents("foo", ExternalAuthzService)
	time.Sleep(time.Duration(50) * time.Second)

	// 2. Deploy the External Authorizer Deployment
	util.Log.Info("Deploy the External Authorizer")
	util.KubeApplyContents("foo", ExternalAuthzDeployment)
	time.Sleep(time.Duration(50) * time.Second)

	// 2.
	util.Log.Info("Verfiy the sample external authorizer is up and running")

	extAuthPod, err := util.GetPodName("foo", "app=ext-authz")
	util.Inspect(err, "Failed to get ext-authz pod name", "", t)
	extAuthPodverf, err := util.Shell("oc logs %s -n foo -c ext-authz", extAuthPod)
	if err != nil {
		t.Fatalf("external autorization is not running: %v", err)
	}
	extAuthPodverf = strings.Trim(extAuthPodverf, "'")

	// ext, err := util.Shell("oc get smcp -n %s -o jsonpath={.items..metadata.name}", meshNamespace)
	// util.Inspect(err, "Failed to smcp name", "", t)
	// util.Shell("oc get configmap istio-%s -n %s", ext, meshNamespace)

	// // oc get configmap istio-basic -o jsonpath={.data} -n istio-system  --> configmap jsonpath
	// 1.Define the external authorizer

	util.Log.Info("Edit configmap for two external providers")
	util.Shell(`kubectl patch -n %s configmap/istio-basic --type merge -p '{"data": {"mesh": "# Add the following content to define the external authorizers.\nextensionProviders:\n- name: \"sample-ext-authz-grpc\"\n  envoyExtAuthzGrpc:\n    service: \"ext-authz.foo.svc.cluster.local\"\n    port: \"9000\"\n- name: \"sample-ext-authz-http\"\n  envoyExtAuthzHttp:\n    service: \"ext-authz.foo.svc.cluster.local\"\n    port: \"8000\"\n    includeRequestHeadersInCheck: [\"x-ext-authz\"]"}}'`, meshNamespace)

	// Enable with external authorization
	// 1.
	util.Log.Info("Enable with external authorization")
	util.KubeApplyContents("foo", ExternalRoute)

	// 2.
	util.Log.Info("Verfiy the header with deny server")
	util.Inspect(err, "Failed to get sleep pod name", "", t)
	deny := fmt.Sprintf(`curl "http://httpbin.foo:8000/headers" -H "x-ext-authz: deny" -s`)
	extAuthDeny, err := util.Shell("kubectl exec %s -c sleep -n foo -- %s", sleepPod, deny)
	if err != nil {
		t.Fatalf("external autorization is not running: %v", err)
	}
	extAuthPodverf = strings.Trim(extAuthDeny, "'")

	// 3.
	util.Log.Info("Verfiy the header with allow server")
	util.Inspect(err, "Failed to get sleep pod name", "", t)
	allow := fmt.Sprintf(`curl "http://httpbin.foo:8000/headers" -H "x-ext-authz: allow" -s`)
	extAuthAllow, err := util.Shell("kubectl exec %s -c sleep -n foo -- %s", sleepPod, allow)
	if err != nil {
		t.Fatalf("external autorization is not running: %v", err)
	}
	extAuthPodverf = strings.Trim(extAuthAllow, "'")

	// 4.
	util.Inspect(err, "Failed to get sleep pod name", "", t)
	cmds := fmt.Sprintf(`curl http://httpbin.foo:8000/ip -sS -o /dev/null -w "%%{http_code}\n"`)
	msges, err := util.PodExec("foo", sleepPod, "sleep", cmds, true)
	util.Inspect(err, "Failed to get response", "", t)
	if !strings.Contains(msges, "200") {
		util.Log.Errorf("Verify setup -- Unexpected response code: %s", msges)
	} else {
		util.Log.Infof("Success. Get expected response: %s", msges)
	}

	//5.

	util.Inspect(err, "Failed to get ext-authz pod name", "", t)
	// util.Shell("kubectl logs %s -n foo -c ext-authz", extAuthPod)

	extAuthPodverfs, err := util.Shell("oc logs %s -n foo -c ext-authz", extAuthPod)
	if err != nil {
		t.Fatalf("external autorization is not running: %v", err)
	}
	extAuthPodverfs = strings.Trim(extAuthPodverfs, "'")

}
