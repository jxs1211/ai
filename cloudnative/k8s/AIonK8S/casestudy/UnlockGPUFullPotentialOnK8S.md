Unlocking the Full Potential of GPUs for AI Workloads on Kubernetes - Kevin Klues, NVIDIA
https://www.youtube.com/watch?v=1QfShSQLsbs

[Device Plugins 2.0: How to Build a Driver for Dynamic Resource Allocation - K Klues & Alexey Fomenko - YouTube](https://www.youtube.com/watch?v=_fi9asserLE)

[NVIDIA/k8s-dra-driver: Dynamic Resource Allocation (DRA) for NVIDIA GPUs in Kubernetes (github.com)](https://github.com/NVIDIA/k8s-dra-driver)

Dynamic Resource Allocation (DRA) is new Kubernetes feature that puts resource scheduling in the hands of 3rd-party developers. It moves away from the limited "countable" interface for requesting access to resources (e.g. "nvidia.com/gpu: 2"), providing an API more akin to that of persistent volumes. 

In the context of GPUs, this unlocks a host of new features without the need for awkward solutions shoehorned on top of the existing device plugin API. These features include: 

- Controlled GPU Sharing (both within a pod and across pods) 
- Multiple GPU models per node (e.g. T4 and A100) 
- Specifying arbitrary constraints for a GPU (min/max memory, device model, etc.) 
- Dynamic allocation of Multi-Instance GPUs (MIG) *

the list goes on ... In this talk, you will learn about the DRA resource driver we have built for GPUs. We walk through each of the features it provides, and conclude with a series of demos showing you how you can get started using it today.



![image-20231212092638003](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20231212092638003.png)

![image-20231212095420426](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20231212095420426.png)

![image-20231218092148515](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20231218092148515.png)