package volumes

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	corev1 "k8s.io/api/core/v1"
)

func TestGetServiceAuthenticationTokenVolume(t *testing.T) {
	expectedExpSeconds := int64(3600)
	expectedVolume := corev1.Volume{
		Name: "projected-sa--foodot-yelpdot-com",
		VolumeSource: corev1.VolumeSource{
			Projected: &corev1.ProjectedVolumeSource{
				Sources: []corev1.VolumeProjection{
					{
						ServiceAccountToken: &corev1.ServiceAccountTokenProjection{
							Audience:          "foo.yelp.com",
							ExpirationSeconds: &expectedExpSeconds,
							Path:              "token",
						},
					},
				},
			},
		},
	}
	expectedMount := corev1.VolumeMount{
		Name:      "projected-sa--foodot-yelpdot-com",
		ReadOnly:  true,
		MountPath: "/var/secret/serviceaccount/foo",
	}
	mockJwtAuthConfig := &sync.Map{}
	mockJwtAuthConfig.Store(
		jwtServiceAuthConfigKey,
		jwtServiceAuthTokenSettings{
			Audience:      "foo.yelp.com",
			ContainerPath: "/var/secret/serviceaccount/foo",
		},
	)
	configReader := &configstore.Store{Data: mockJwtAuthConfig}
	outputMount, outputVolume, _ := GetServiceAuthenticationTokenVolume(configReader)
	if !reflect.DeepEqual(outputMount, expectedMount) {
		t.Errorf("Wrong SA volume mount.\nExpected: %+v\nGot: %+v", expectedMount, outputMount)
	}
	if !reflect.DeepEqual(outputVolume, expectedVolume) {
		t.Errorf("Wrong SA volume config.\nExpected: %+v\nGot: %+v", expectedVolume, outputVolume)
	}
}

func TestFormatServiceAccountVolumeName(t *testing.T) {
	testCases := []string{
		"foobar",
		"foo/yelp/com",
		"/foobar_yelp",
		"foo.yelp.com",
		"foooooooooooooooooooooooooooooooobaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaar",
	}
	expected := []string{
		"projected-sa--foobar",
		"projected-sa--foo-yelp-com",
		"projected-sa--foobar--yelp",
		"projected-sa--foodot-yelpdot-com",
		"projected-sa--foooooooooooooooooooooooooooooooobaaaaaaaaaaaaaaa",
	}
	for i, testCase := range testCases {
		if output := formatServiceAccountVolumeName(testCase); output != expected[i] {
			t.Errorf("Wrong SA volume name format.\nExpected: %s\nGot: %s", expected[i], output)
		}
	}
}

func TestServiceRequiresAuthenticationToken(t *testing.T) {
	authenticatingServicesCache = map[string]bool{"foobar": true}
	authenticatingServicesLastLoaded = time.Now()
	testCases := []string{
		"foobar",
		"bizzbuzz",
	}
	expected := []bool{
		true,
		false,
	}
	for i, testCase := range testCases {
		if output, _ := ServiceRequiresAuthenticationToken(testCase); output != expected[i] {
			t.Errorf("Error checking service %s.\nExpected: %v\nGot: %v", testCase, expected[i], output)
		}
	}
}

func TestGetAuthenticatingServices(t *testing.T) {
	authenticatingServicesLastLoaded = time.Time{} // invalidate cache
	expected := map[string]bool{"foobar": true, "bizzbuzz": true}
	tempConf, err := os.CreateTemp("", "paasta-tools-tests")
	if err != nil {
		t.Fatalf("Error creating temp test config: %s", err)
	}
	tempConfName := tempConf.Name()
	defer os.Remove(tempConfName)
	tempConf.WriteString("services:\n- foobar\n- bizzbuzz\n")
	tempConf.Close()
	result, err := getAuthenticatingServices(tempConfName)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Config not properly loaded.\nExpected: %+v\nGot: %+v", expected, result)
	}
	if !reflect.DeepEqual(result, authenticatingServicesCache) {
		t.Errorf("Config values not properly cached.\nExpected: %+v\nGot: %+v", result, authenticatingServicesCache)
	}
}
