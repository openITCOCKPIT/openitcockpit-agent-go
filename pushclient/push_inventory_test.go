package pushclient

import (
	"context"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/packagemanager"
)

func TestPushPackageInfoKeepsInstalledKBs(t *testing.T) {
	fake := &fakePushServer{}
	server := httptest.NewServer(fake)
	defer server.Close()

	target, _ := url.Parse(server.URL + "/agentconnector/updateSoftware.json")
	p := &PushClient{
		authConfiguration:    authConfiguration{Password: "secret", UUID: "test"},
		urlSubmitPackageInfo: target,
		timeout:              5 * time.Second,
	}

	p.pushPackageInfo(context.Background(), packagemanager.PackageInfo{WindowsInstalledKBs: []string{"5079473"}})

	if len(fake.payloads) != 1 {
		t.Fatalf("expected one request, got %d", len(fake.payloads))
	}
	checkData, ok := fake.payloads[0]["checkdata"].(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected payload: %v", fake.payloads[0])
	}
	expected := []interface{}{"5079473"}
	if !reflect.DeepEqual(checkData["WindowsInstalledKBs"], expected) {
		t.Errorf("expected %v, got %v", expected, checkData["WindowsInstalledKBs"])
	}
}
