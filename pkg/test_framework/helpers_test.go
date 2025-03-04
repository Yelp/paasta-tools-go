package framework

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	fizzbuzzYaml = `
apiVersion: fizz/v1
kind: buzz
metadata:
  annotations:
    fizz.foo: '0'
    fizz.bar: "false"
  labels:
    buzz.hello: world
  name: fizzbuzz
  namespace: helloworld
spec:
  trigger:
    - fizz
    - buzz
    - fizzbuzz
  fizz: 3
  buzz: 5
  fizzbuzz: "15"
`
	fizzbuzzJson = `
{
  "apiVersion": "fizz/v1",
  "kind": "buzz",
  "metadata": {
    "annotations": {
      "fizz.foo": "0",
      "fizz.bar": "false"
    },
    "labels": {
      "buzz.hello": "world"
    },
    "name": "fizzbuzz",
    "namespace": "helloworld"
  },
  "spec": {
    "trigger": [
      "fizz",
      "buzz",
      "fizzbuzz"
    ],
    "fizz": 3,
    "buzz": 5,
    "fizzbuzz": "15"
  }
}
`

	fizzbuzzMini = `{` +
		`"apiVersion":"fizz/v1",` +
		`"kind":"buzz",` +
		`"metadata":{` +
		`"annotations":{"fizz.bar":"false","fizz.foo":"0"},` +
		`"labels":{"buzz.hello":"world"},` +
		`"name":"fizzbuzz",` +
		`"namespace":"helloworld"` +
		`},` +
		`"spec":{` +
		`"buzz":5,"fizz":3,"fizzbuzz":"15","trigger":["fizz","buzz","fizzbuzz"]` +
		`}}
`

	fizzbuzzSrvYaml = `
apiVersion: v1
kind: Service
metadata:
  labels:
    buzz.hello: world
  name: fizzbuzz
  namespace: helloworld
spec:
  ports:
  - name: rest
    port: 8080
    protocol: TCP
    targetPort: 8080
`

	fizzbuzzSrvMini = `{` +
		`"apiVersion":"v1",` +
		`"kind":"Service",` +
		`"metadata":{` +
		`"creationTimestamp":null,` +
		`"labels":{"buzz.hello":"world"},` +
		`"name":"fizzbuzz",` +
		`"namespace":"helloworld"` +
		`},` +
		`"spec":{` +
		`"ports":[{"name":"rest","port":8080,"protocol":"TCP","targetPort":8080}]` +
		`},` +
		`"status":{"loadBalancer":{}}}
`
)

func TestLoadUnstructured(t *testing.T) {
	reader := bytes.NewReader([]byte(fizzbuzzYaml))
	obj1, err := LoadUnstructured(reader)
	assert.NoError(t, err)
	assert.Equal(t, "buzz", obj1.GetKind())
	assert.Equal(t, "fizz/v1", obj1.GetAPIVersion())
	assert.Equal(t, "fizzbuzz", obj1.GetName())
	assert.Equal(t, "helloworld", obj1.GetNamespace())
	assert.Equal(t, map[string]string{
		"buzz.hello": "world",
	}, obj1.GetLabels())
	assert.Equal(t, map[string]string{
		"fizz.foo": "0",
		"fizz.bar": "false",
	}, obj1.GetAnnotations())
	json1, _ := obj1.MarshalJSON()
	assert.Equal(t, fizzbuzzMini, string(json1))

	reader.Reset([]byte(fizzbuzzJson))
	obj2, err := LoadUnstructured(reader)
	assert.NoError(t, err)
	json2, _ := obj2.MarshalJSON()
	assert.Equal(t, fizzbuzzMini, string(json2))

	reader.Reset([]byte(fizzbuzzMini))
	obj3, err := LoadUnstructured(reader)
	assert.NoError(t, err)
	json3, _ := obj3.MarshalJSON()
	assert.Equal(t, fizzbuzzMini, string(json3))

	reader.Reset([]byte(fizzbuzzYaml))
	obj4 := &unstructured.Unstructured{}
	err = LoadInto(reader, obj4)
	assert.NoError(t, err)
	json4, _ := obj4.MarshalJSON()
	assert.Equal(t, fizzbuzzMini, string(json4))

	reader.Reset([]byte(fizzbuzzJson))
	obj5 := &unstructured.Unstructured{}
	err = LoadInto(reader, obj5)
	assert.NoError(t, err)
	json5, _ := obj5.MarshalJSON()
	assert.Equal(t, fizzbuzzMini, string(json5))
}

// Only test Service because it is impractical to test loading all types of objects
// - we cede the responsibility for that to sigs.k8s.io/yaml or whatever kubernetes
// is using internally
func TestLoadService(t *testing.T) {
	reader := bytes.NewReader([]byte(fizzbuzzSrvYaml))
	service := &corev1.Service{}
	err := LoadInto(reader, service)
	assert.NoError(t, err)

	// rather than json.Marshall, see if ToUnstructured conversion works for our object
	tmp, err := runtime.DefaultUnstructuredConverter.ToUnstructured(service)
	assert.NoError(t, err)
	obj1 := unstructured.Unstructured{Object: tmp}
	assert.Equal(t, "Service", obj1.GetKind())
	assert.Equal(t, "v1", obj1.GetAPIVersion())
	assert.Equal(t, "fizzbuzz", obj1.GetName())
	assert.Equal(t, "helloworld", obj1.GetNamespace())
	json2, _ := obj1.MarshalJSON()
	assert.Equal(t, fizzbuzzSrvMini, string(json2))
}

func TestReadWriteDeleteValue(t *testing.T) {
	reader := bytes.NewReader([]byte(fizzbuzzSrvYaml))
	service := &unstructured.Unstructured{}
	err := LoadInto(reader, service)
	assert.NoError(t, err)

	label, err := ReadValue(service, "metadata", "labels", "buzz.hello")
	assert.NoError(t, err)
	assert.Equal(t, "world", label)

	// Read a value nested under value
	_, err = ReadValue(service, "metadata", "labels", "buzz.hello", "dummy")
	assert.NotNil(t, err)

	// Read non-existing value
	_, err = ReadValue(service, "metadata", "labels", "dummy")
	assert.NotNil(t, err)

	// Read a value under non-existing section
	_, err = ReadValue(service, "metadata", "annotations", "dummy")
	assert.NotNil(t, err)

	// Delete non-existing section - not an error
	err = DeleteValue(service, "metadata", "annotations", "dummy")
	assert.NoError(t, err)

	// Delete non-existing value - not an error
	err = DeleteValue(service, "metadata", "labels", "dummy")
	assert.NoError(t, err)

	// Write a value
	err = WriteValue(service, "value", "metadata", "labels", "dummy")
	assert.NoError(t, err)

	// Check it is there
	tmp, err := ReadValue(service, "metadata", "labels", "dummy")
	assert.NoError(t, err)
	assert.Equal(t, "value", tmp)

	// Overwrite a value
	err = WriteValue(service, "bar", "metadata", "labels", "dummy")
	assert.NoError(t, err)

	// Read whole section
	labels, err := ReadValue(service, "metadata", "labels")
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"buzz.hello": "world",
		"dummy":      "bar",
	}, labels)

	// Delete dummy value now
	err = DeleteValue(service, "metadata", "labels", "dummy")
	assert.NoError(t, err)

	// Read whole section again
	labels, err = ReadValue(service, "metadata", "labels")
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"buzz.hello": "world",
	}, labels)

	// Delete whole section
	err = DeleteValue(service, "metadata", "labels")
	assert.NoError(t, err)

	// Read non-existing section
	labels, err = ReadValue(service, "metadata", "labels")
	assert.NotNil(t, err)

	// Write whole section
	labels = map[string]interface{}{
		"fizz.hello": "there",
	}
	err = WriteValue(service, labels, "metadata", "labels")
	assert.NoError(t, err)

	// Read whole section again
	labels, err = ReadValue(service, "metadata", "labels")
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"fizz.hello": "there",
	}, labels)

	// Read complex value
	ports, err := ReadValue(service, "spec", "ports")
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"name":       "rest",
		"port":       int64(8080),
		"protocol":   "TCP",
		"targetPort": int64(8080),
	}, ports.([]interface{})[0])

	// Write complex value
	ports.([]interface{})[0].(map[string]interface{})["name"] = "blob"
	err = WriteValue(service, ports, "spec", "ports")
	assert.NoError(t, err)

	// Verify the new value
	ports, err = ReadValue(service, "spec", "ports")
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"name":       "blob",
		"port":       int64(8080),
		"protocol":   "TCP",
		"targetPort": int64(8080),
	}, ports.([]interface{})[0])
}
