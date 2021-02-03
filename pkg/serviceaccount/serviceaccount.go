package serviceaccount

import (
	"context"
	"regexp"

	"github.com/Yelp/paasta-tools-go/pkg/hashing"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var serviceAccountRegex = regexp.MustCompile("[^0-9a-zA-Z]+")

func EnsureServiceAccountForIamRole(ctx context.Context, iam_role string, namespace string, client controllerruntimeclient.Client) error {
	serviceAccountName := getServiceAccountNameForIamRole(iam_role)

	// check if service account with this name+namespace already exists
	glog.V(4).Infof("Looking for service account called %s in namespace %s", serviceAccountName, namespace)
	result := &corev1.ServiceAccount{}
	err := client.Get(ctx,
		types.NamespacedName{
			Name:      serviceAccountName,
			Namespace: namespace,
		},
		result,
	)
	if err != nil {
		glog.Errorf("%s", err)
		if errors.IsNotFound(err) {
			glog.Infof("Service account not found, creating it\n")
			annotations := map[string]string{
				"eks.amazonaws.com/role-arn": iam_role,
			}
			sa := &corev1.ServiceAccount{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name:        serviceAccountName,
					Namespace:   namespace,
					Annotations: annotations,
				},
				Secrets:                      nil,
				ImagePullSecrets:             nil,
				AutomountServiceAccountToken: nil,
			}
			objectHash, err := hashing.ComputeHashForKubernetesObject(sa)
			if err != nil {
				return err
			}
			hashing.SetKubernetesObjectHash(objectHash, sa)
			err = client.Create(ctx, sa)
			if err != nil {
				return err
			}
			return nil
		}
		// some other error occurred
		return err
	}

	// service account already exists so do nothing
	glog.V(4).Info("Service account already exists")
	return nil
}

func getServiceAccountNameForIamRole(iamRole string) (serviceAccountName string) {
	serviceAccountName = serviceAccountRegex.ReplaceAllString(iamRole, "-")
	return serviceAccountName
}
