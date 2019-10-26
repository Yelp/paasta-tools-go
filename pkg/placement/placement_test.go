package placement

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestSetPlacementConstraintsUniqueNodes(test *testing.T) {
	service := "foo"
	instance := "bar"
	expected := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					corev1.WeightedPodAffinityTerm{
						Weight: 1000,
						PodAffinityTerm: corev1.PodAffinityTerm{
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									metav1.LabelSelectorRequirement{
										Key:      "paasta.yelp.com/service",
										Operator: "In",
										Values:   []string{service},
									},
									metav1.LabelSelectorRequirement{
										Key:      "paasta.yelp.com/instance",
										Operator: "In",
										Values:   []string{instance},
									},
								},
							},
							TopologyKey: "yelp.com/hostname",
						},
					},
				},
			},
		},
	}
	actual := corev1.PodSpec{}
	SetPlacementConstraints(&actual, UniqueNodes, service, instance)
	if !reflect.DeepEqual(expected, actual) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", expected, actual)
	}
}

func TestSetPlacementConstraintsRequireUniqueNodes(test *testing.T) {
	service := "foo"
	instance := "bar"
	expected := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					corev1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								metav1.LabelSelectorRequirement{
									Key:      "paasta.yelp.com/service",
									Operator: "In",
									Values:   []string{service},
								},
								metav1.LabelSelectorRequirement{
									Key:      "paasta.yelp.com/instance",
									Operator: "In",
									Values:   []string{instance},
								},
							},
						},
						TopologyKey: "yelp.com/hostname",
					},
				},
			},
		},
	}
	actual := corev1.PodSpec{}
	SetPlacementConstraints(&actual, RequireUniqueNodes, service, instance)
	if !reflect.DeepEqual(expected, actual) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", expected, actual)
	}
}

func TestSetPlacementConstraintsUniqueAvailabilityZones(test *testing.T) {
	service := "foo"
	instance := "bar"
	expected := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					corev1.WeightedPodAffinityTerm{
						Weight: 1000,
						PodAffinityTerm: corev1.PodAffinityTerm{
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									metav1.LabelSelectorRequirement{
										Key:      "paasta.yelp.com/service",
										Operator: "In",
										Values:   []string{service},
									},
									metav1.LabelSelectorRequirement{
										Key:      "paasta.yelp.com/instance",
										Operator: "In",
										Values:   []string{instance},
									},
								},
							},
							TopologyKey: "yelp.com/habitat",
						},
					},
					corev1.WeightedPodAffinityTerm{
						Weight: 1000,
						PodAffinityTerm: corev1.PodAffinityTerm{
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									metav1.LabelSelectorRequirement{
										Key:      "paasta.yelp.com/service",
										Operator: "In",
										Values:   []string{service},
									},
									metav1.LabelSelectorRequirement{
										Key:      "paasta.yelp.com/instance",
										Operator: "In",
										Values:   []string{instance},
									},
								},
							},
							TopologyKey: "yelp.com/hostname",
						},
					},
				},
			},
		},
	}
	actual := corev1.PodSpec{}
	SetPlacementConstraints(&actual, UniqueAvailabilityZones, service, instance)
	if !reflect.DeepEqual(expected, actual) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", expected, actual)
	}
}

func TestSetPlacementConstraintsRequireUniqueAvailabilityZones(test *testing.T) {
	service := "foo"
	instance := "bar"
	expected := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					corev1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								metav1.LabelSelectorRequirement{
									Key:      "paasta.yelp.com/service",
									Operator: "In",
									Values:   []string{service},
								},
								metav1.LabelSelectorRequirement{
									Key:      "paasta.yelp.com/instance",
									Operator: "In",
									Values:   []string{instance},
								},
							},
						},
						TopologyKey: "yelp.com/habitat",
					},
				},
			},
		},
	}
	actual := corev1.PodSpec{}
	SetPlacementConstraints(&actual, RequireUniqueAvailabilityZones, service, instance)
	if !reflect.DeepEqual(expected, actual) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", expected, actual)
	}
}
