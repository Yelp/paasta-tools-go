package placement

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PlacementConstraint int

const (
	None PlacementConstraint = iota
	UniqueNodes
	UniqueAvailabilityZones
	RequireUniqueNodes
	RequireUniqueAvailabilityZones
)

type podAffinityTermsSetter func(*corev1.PodAntiAffinity, *[]corev1.PodAffinityTerm)

const (
	podAffinityWeight int32 = 1000
)

func setPlacementConstraints(spec *corev1.PodSpec, setter podAffinityTermsSetter, service string, instance string, topologyKeys ...string) {
	podAffinityTerms := []corev1.PodAffinityTerm{}
	for _, topologyKey := range topologyKeys {
		podAffinityTerms = append(
			podAffinityTerms,
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
				TopologyKey: topologyKey,
			})
	}

	if spec.Affinity == nil {
		spec.Affinity = &corev1.Affinity{}
	}
	affinity := spec.Affinity
	if affinity.PodAntiAffinity == nil {
		affinity.PodAntiAffinity = &corev1.PodAntiAffinity{}
	}
	setter(affinity.PodAntiAffinity, &podAffinityTerms)
}

func setRequired(podAntiAffinity *corev1.PodAntiAffinity, podAffinityTerms *[]corev1.PodAffinityTerm) {
	if podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []corev1.PodAffinityTerm{}
	}
	podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
		podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution,
		*podAffinityTerms...)
}

func setPreferred(podAntiAffinity *corev1.PodAntiAffinity, podAffinityTerms *[]corev1.PodAffinityTerm) {
	if podAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution == nil {
		podAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.WeightedPodAffinityTerm{}
	}
	for _, podAffinityTerm := range *podAffinityTerms {
		podAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
			podAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution,
			corev1.WeightedPodAffinityTerm{
				PodAffinityTerm: podAffinityTerm,
				Weight:          podAffinityWeight,
			})
	}
}

func SetPlacementConstraints(spec *corev1.PodSpec, placementConstraint PlacementConstraint, service string, instance string) {
	switch placementConstraint {
	case None:
	case UniqueNodes:
		setPlacementConstraints(spec, setPreferred, service, instance, "yelp.com/hostname")
	case UniqueAvailabilityZones:
		// At least try to run on different nodes if the unique habitats preference cannot be satisfied.
		setPlacementConstraints(spec, setPreferred, service, instance, "yelp.com/habitat", "yelp.com/hostname")
	case RequireUniqueNodes:
		setPlacementConstraints(spec, setRequired, service, instance, "yelp.com/hostname")
	case RequireUniqueAvailabilityZones:
		setPlacementConstraints(spec, setRequired, service, instance, "yelp.com/habitat")
	}
}
