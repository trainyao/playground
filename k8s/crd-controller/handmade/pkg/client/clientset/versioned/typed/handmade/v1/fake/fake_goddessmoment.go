// this is handmade crd & controller

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	handmadev1 "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/apis/handmade/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeGoddessMoments implements GoddessMomentInterface
type FakeGoddessMoments struct {
	Fake *FakeHandmadeV1
	ns   string
}

var goddessmomentsResource = schema.GroupVersionResource{Group: "handmade.crd.playground.trainyao.io", Version: "v1", Resource: "goddessmoments"}

var goddessmomentsKind = schema.GroupVersionKind{Group: "handmade.crd.playground.trainyao.io", Version: "v1", Kind: "GoddessMoment"}

// Get takes name of the goddessMoment, and returns the corresponding goddessMoment object, and an error if there is any.
func (c *FakeGoddessMoments) Get(ctx context.Context, name string, options v1.GetOptions) (result *handmadev1.GoddessMoment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(goddessmomentsResource, c.ns, name), &handmadev1.GoddessMoment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*handmadev1.GoddessMoment), err
}

// List takes label and field selectors, and returns the list of GoddessMoments that match those selectors.
func (c *FakeGoddessMoments) List(ctx context.Context, opts v1.ListOptions) (result *handmadev1.GoddessMomentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(goddessmomentsResource, goddessmomentsKind, c.ns, opts), &handmadev1.GoddessMomentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &handmadev1.GoddessMomentList{ListMeta: obj.(*handmadev1.GoddessMomentList).ListMeta}
	for _, item := range obj.(*handmadev1.GoddessMomentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested goddessMoments.
func (c *FakeGoddessMoments) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(goddessmomentsResource, c.ns, opts))

}

// Create takes the representation of a goddessMoment and creates it.  Returns the server's representation of the goddessMoment, and an error, if there is any.
func (c *FakeGoddessMoments) Create(ctx context.Context, goddessMoment *handmadev1.GoddessMoment, opts v1.CreateOptions) (result *handmadev1.GoddessMoment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(goddessmomentsResource, c.ns, goddessMoment), &handmadev1.GoddessMoment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*handmadev1.GoddessMoment), err
}

// Update takes the representation of a goddessMoment and updates it. Returns the server's representation of the goddessMoment, and an error, if there is any.
func (c *FakeGoddessMoments) Update(ctx context.Context, goddessMoment *handmadev1.GoddessMoment, opts v1.UpdateOptions) (result *handmadev1.GoddessMoment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(goddessmomentsResource, c.ns, goddessMoment), &handmadev1.GoddessMoment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*handmadev1.GoddessMoment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeGoddessMoments) UpdateStatus(ctx context.Context, goddessMoment *handmadev1.GoddessMoment, opts v1.UpdateOptions) (*handmadev1.GoddessMoment, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(goddessmomentsResource, "status", c.ns, goddessMoment), &handmadev1.GoddessMoment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*handmadev1.GoddessMoment), err
}

// Delete takes name of the goddessMoment and deletes it. Returns an error if one occurs.
func (c *FakeGoddessMoments) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(goddessmomentsResource, c.ns, name, opts), &handmadev1.GoddessMoment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGoddessMoments) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(goddessmomentsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &handmadev1.GoddessMomentList{})
	return err
}

// Patch applies the patch and returns the patched goddessMoment.
func (c *FakeGoddessMoments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *handmadev1.GoddessMoment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(goddessmomentsResource, c.ns, name, pt, data, subresources...), &handmadev1.GoddessMoment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*handmadev1.GoddessMoment), err
}
