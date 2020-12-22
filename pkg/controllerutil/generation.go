package controllerutil

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var SetControllerReference = controllerutil.SetControllerReference

var (
	AnnotationControllerGeneration = "controller-generation"
	AnnotationRestartedAt          = "kubectl.kubernetes.io/restartedAt"
)

func IsControllerGenerationEqual(cur metav1.Object, next metav1.Object) bool {
	if nextOwner := metav1.GetControllerOf(next); nextOwner != nil {
		if curOwner := metav1.GetControllerOf(cur); curOwner != nil {
			if curOwner.UID != nextOwner.UID {
				return false
			}
		}
	}

	annotations := cur.GetAnnotations()
	nextAnnotations := next.GetAnnotations()

	return isEqualProp(annotations, nextAnnotations, AnnotationControllerGeneration) && isEqualProp(annotations, nextAnnotations, AnnotationRestartedAt)
}

func isEqualProp(cur map[string]string, next map[string]string, prop string) bool {
	if cur == nil {
		return false
	}
	if next == nil {
		return false
	}
	return cur[prop] == next[prop]
}
