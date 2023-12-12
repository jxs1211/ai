Device Plugins 2.0: How to Build a Driver for Dynamic Resource Allocation - K Klues & Alexey Fomenko
https://www.youtube.com/watch?v=_fi9asserLE

![image-20231213092753732](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20231213092753732.png)



[enhancements/keps/sig-node/3063-dynamic-resource-allocation at master · kubernetes/enhancements (github.com)](https://github.com/kubernetes/enhancements/tree/master/keps/sig-node/3063-dynamic-resource-allocation)

[cncf-tags/container-device-interface (github.com)](https://github.com/cncf-tags/container-device-interface)

[NVIDIA/k8s-dra-driver: Dynamic Resource Allocation (DRA) for NVIDIA GPUs in Kubernetes (github.com)](https://github.com/NVIDIA/k8s-dra-driver)

staging/src/k8s.io/dynamic-resource-allocation/controller/controller.go

```go

// Driver provides the actual allocation and deallocation operations.
type Driver interface {
	// GetClassParameters gets called to retrieve the parameter object
	// referenced by a class. The content should be validated now if
	// possible. class.Parameters may be nil.
	//
	// The caller will wrap the error to include the parameter reference.
	GetClassParameters(ctx context.Context, class *resourcev1alpha1.ResourceClass) (interface{}, error)

	// GetClaimParameters gets called to retrieve the parameter object
	// referenced by a claim. The content should be validated now if
	// possible. claim.Spec.Parameters may be nil.
	//
	// The caller will wrap the error to include the parameter reference.
	GetClaimParameters(ctx context.Context, claim *resourcev1alpha1.ResourceClaim, class *resourcev1alpha1.ResourceClass, classParameters interface{}) (interface{}, error)

	// Allocate gets called when a ResourceClaim is ready to be allocated.
	// The selectedNode is empty for ResourceClaims with immediate
	// allocation, in which case the resource driver decides itself where
	// to allocate. If there is already an on-going allocation, the driver
	// may finish it and ignore the new parameters or abort the on-going
	// allocation and try again with the new parameters.
	//
	// Parameters have been retrieved earlier.
	//
	// If selectedNode is set, the driver must attempt to allocate for that
	// node. If that is not possible, it must return an error. The
	// controller will call UnsuitableNodes and pass the new information to
	// the scheduler, which then will lead to selecting a diffent node
	// if the current one is not suitable.
	//
	// The objects are read-only and must not be modified. This call
	// must be idempotent.
	Allocate(ctx context.Context, claim *resourcev1alpha1.ResourceClaim, claimParameters interface{}, class *resourcev1alpha1.ResourceClass, classParameters interface{}, selectedNode string) (*resourcev1alpha1.AllocationResult, error)

	// Deallocate gets called when a ResourceClaim is ready to be
	// freed.
	//
	// The claim is read-only and must not be modified. This call must be
	// idempotent. In particular it must not return an error when the claim
	// is currently not allocated.
	//
	// Deallocate may get called when a previous allocation got
	// interrupted. Deallocate must then stop any on-going allocation
	// activity and free resources before returning without an error.
	Deallocate(ctx context.Context, claim *resourcev1alpha1.ResourceClaim) error

	// UnsuitableNodes checks all pending claims with delayed allocation
	// for a pod. All claims are ready for allocation by the driver
	// and parameters have been retrieved.
	//
	// The driver may consider each claim in isolation, but it's better
	// to mark nodes as unsuitable for all claims if it not all claims
	// can be allocated for it (for example, two GPUs requested but
	// the node only has one).
	//
	// The result of the check is in ClaimAllocation.UnsuitableNodes.
	// An error indicates that the entire check must be repeated.
	UnsuitableNodes(ctx context.Context, pod *v1.Pod, claims []*ClaimAllocation, potentialNodes []string) error
}
```



