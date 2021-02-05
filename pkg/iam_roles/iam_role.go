package iam_role

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"regexp"
	controllerruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// IamRoleConfig : config for AWS IAM role settings for a paasta container
type IamRoleConfig struct {
	// +optional
	IamRoleProvider *string `json:"iam_role_provider,omitempty"`
	// +optional
	IamRole *string `json:"iam_role,omitempty"`
	// PAASTA-16919: remove everything related to fs_group when
	// https://github.com/aws/amazon-eks-pod-identity-webhook/issues/8 will be
	// fixed.
	// +optional
	FsGroup *int64 `json:"fs_group,omitempty"`
}

var defaultIamRoleProvider = "kiam"
var DefaultFsGroup int64 = 65534

var serviceAccountRegex = regexp.MustCompile("[^0-9a-zA-Z]+")

// SetIamRoleConfigDefaults: sets the default values for the AWS IAM role config
func SetIamRoleConfigDefaults(iamRoleConfig *IamRoleConfig) {
	if iamRoleConfig.IamRoleProvider == nil {
		iamRoleConfig.IamRoleProvider = &defaultIamRoleProvider
	}
	if iamRoleConfig.FsGroup == nil {
		iamRoleConfig.FsGroup = &DefaultFsGroup
	}
}

// EnsureForIamRole: prepare AWS IAM role for use
func EnsureForIamRole(ctx context.Context, client controllerruntimeclient.Client, namespace string, iamRoleConfig *IamRoleConfig) error {
	// We need to create a service account only for "aws" provider
	if iamRoleConfig.IamRoleProvider != nil && *iamRoleConfig.IamRoleProvider == "aws" {
		glog.V(4).Infof("Ensuring service account in %s for iam_role %v exists", namespace, iamRoleConfig.IamRole)
		if iamRoleConfig.IamRole == nil {
			return fmt.Errorf("%s/%v: iam_role must be specified when iam_role_provider is set to 'aws'", namespace, iamRoleConfig.IamRole)
		}

		saName := getServiceAccountNameForIamRole(iamRoleConfig.IamRole)

		// check if service account with this name+namespace already exists
		glog.V(4).Infof("Looking for service account called %s in namespace %s", saName, namespace)
		result := &corev1.ServiceAccount{}
		err := client.Get(ctx,
			types.NamespacedName{
				Name:      saName,
				Namespace: namespace,
			},
			result,
		)

		if err != nil {
			glog.Errorf("Error while getting service account: %s", err)
			if errors.IsNotFound(err) {
				glog.Infof("Service account not found, creating it")
				annotations := map[string]string{
					"eks.amazonaws.com/role-arn": *iamRoleConfig.IamRole,
				}
				service_account := &corev1.ServiceAccount{
					TypeMeta: v1.TypeMeta{},
					ObjectMeta: v1.ObjectMeta{
						Name:        saName,
						Namespace:   namespace,
						Annotations: annotations,
					},
					Secrets:                      nil,
					ImagePullSecrets:             nil,
					AutomountServiceAccountToken: nil,
				}
				err = client.Create(ctx, service_account)
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
	}
	return nil

}

// UpdatePodTemplateSpecForIamRole: updates provided pod template specs for kiam or AWS pod identity
func UpdatePodTemplateSpecForIamRole(podTemplateSpec *corev1.PodTemplateSpec, iamRoleConfig *IamRoleConfig, defaultIamRole string) {
	if iamRoleConfig.IamRoleProvider != nil && *iamRoleConfig.IamRoleProvider == "aws" {
		var fsGroup *int64
		if iamRoleConfig.FsGroup != nil {
			fsGroup = iamRoleConfig.FsGroup
		} else {
			fsGroup = &DefaultFsGroup
		}
		podTemplateSpec.Spec.SecurityContext = &corev1.PodSecurityContext{FSGroup: fsGroup}

		// generate "normalized" SA name from iamRole
		podTemplateSpec.Spec.ServiceAccountName = getServiceAccountNameForIamRole(iamRoleConfig.IamRole)
	} else {
		var iamRole *string
		if iamRoleConfig.IamRole != nil {
			iamRole = iamRoleConfig.IamRole
		} else {
			iamRole = &defaultIamRole
		}
		if podTemplateSpec.Annotations == nil {
			podTemplateSpec.Annotations = map[string]string{}
		}
		podTemplateSpec.Annotations["iam.amazonaws.com/role"] = *iamRole
		podTemplateSpec.Spec.SecurityContext = &corev1.PodSecurityContext{}
		podTemplateSpec.Spec.ServiceAccountName = ""
	}
	return
}

func getServiceAccountNameForIamRole(iamRole *string) (serviceAccountName string) {
	serviceAccountName = serviceAccountRegex.ReplaceAllString(*iamRole, "-")
	return "paasta--" + serviceAccountName
}
