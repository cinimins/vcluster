// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/loft-sh/api/v3/pkg/apis/storage/v1"
	scheme "github.com/loft-sh/api/v3/pkg/client/clientset_generated/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ClusterAccessesGetter has a method to return a ClusterAccessInterface.
// A group's client should implement this interface.
type ClusterAccessesGetter interface {
	ClusterAccesses() ClusterAccessInterface
}

// ClusterAccessInterface has methods to work with ClusterAccess resources.
type ClusterAccessInterface interface {
	Create(ctx context.Context, clusterAccess *v1.ClusterAccess, opts metav1.CreateOptions) (*v1.ClusterAccess, error)
	Update(ctx context.Context, clusterAccess *v1.ClusterAccess, opts metav1.UpdateOptions) (*v1.ClusterAccess, error)
	UpdateStatus(ctx context.Context, clusterAccess *v1.ClusterAccess, opts metav1.UpdateOptions) (*v1.ClusterAccess, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ClusterAccess, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ClusterAccessList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ClusterAccess, err error)
	ClusterAccessExpansion
}

// clusterAccesses implements ClusterAccessInterface
type clusterAccesses struct {
	client rest.Interface
}

// newClusterAccesses returns a ClusterAccesses
func newClusterAccesses(c *StorageV1Client) *clusterAccesses {
	return &clusterAccesses{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterAccess, and returns the corresponding clusterAccess object, and an error if there is any.
func (c *clusterAccesses) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ClusterAccess, err error) {
	result = &v1.ClusterAccess{}
	err = c.client.Get().
		Resource("clusteraccesses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterAccesses that match those selectors.
func (c *clusterAccesses) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ClusterAccessList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ClusterAccessList{}
	err = c.client.Get().
		Resource("clusteraccesses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterAccesses.
func (c *clusterAccesses) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("clusteraccesses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a clusterAccess and creates it.  Returns the server's representation of the clusterAccess, and an error, if there is any.
func (c *clusterAccesses) Create(ctx context.Context, clusterAccess *v1.ClusterAccess, opts metav1.CreateOptions) (result *v1.ClusterAccess, err error) {
	result = &v1.ClusterAccess{}
	err = c.client.Post().
		Resource("clusteraccesses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterAccess).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a clusterAccess and updates it. Returns the server's representation of the clusterAccess, and an error, if there is any.
func (c *clusterAccesses) Update(ctx context.Context, clusterAccess *v1.ClusterAccess, opts metav1.UpdateOptions) (result *v1.ClusterAccess, err error) {
	result = &v1.ClusterAccess{}
	err = c.client.Put().
		Resource("clusteraccesses").
		Name(clusterAccess.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterAccess).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *clusterAccesses) UpdateStatus(ctx context.Context, clusterAccess *v1.ClusterAccess, opts metav1.UpdateOptions) (result *v1.ClusterAccess, err error) {
	result = &v1.ClusterAccess{}
	err = c.client.Put().
		Resource("clusteraccesses").
		Name(clusterAccess.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterAccess).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the clusterAccess and deletes it. Returns an error if one occurs.
func (c *clusterAccesses) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clusteraccesses").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterAccesses) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("clusteraccesses").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched clusterAccess.
func (c *clusterAccesses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ClusterAccess, err error) {
	result = &v1.ClusterAccess{}
	err = c.client.Patch(pt).
		Resource("clusteraccesses").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
