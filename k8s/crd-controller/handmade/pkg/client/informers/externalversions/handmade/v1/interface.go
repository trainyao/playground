// this is handmade crd & controller

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// GoddessMoments returns a GoddessMomentInformer.
	GoddessMoments() GoddessMomentInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// GoddessMoments returns a GoddessMomentInformer.
func (v *version) GoddessMoments() GoddessMomentInformer {
	return &goddessMomentInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}