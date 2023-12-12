https://itnext.io/k8sgpt-localai-unlock-kubernetes-superpowers-for-free-584790de9b65

# K8sGPT + LocalAI: Unlock Kubernetes superpowers for free!

[![Tyler](https://miro.medium.com/v2/da:true/resize:fill:88:88/0*dLPyC4_L9JJ3b2Gd)](https://medium.com/@tyler_97636?source=post_page-----584790de9b65--------------------------------)[![ITNEXT](https://miro.medium.com/v2/resize:fill:48:48/1*yAqDFIFA5F_NXalOJKz4TA.png)](https://itnext.io/?source=post_page-----584790de9b65--------------------------------)

[Tyler](https://medium.com/@tyler_97636?source=post_page-----584790de9b65--------------------------------)·Follow

Published in[ITNEXT](https://itnext.io/?source=post_page-----584790de9b65--------------------------------)·8 min read·Apr 27, 2023



321

2







As we all know, LLMs are trending like crazy and the hype is not unjustified. Tons of cool projects leveraging LLM-based text generation are emerging by the day — in fact, I wouldn’t be surprised if another awesome new tool was published during the time it took me to write this blog :)

For the unbeliever, I say that the hype is justified because these projects are not just gimmicks. They are unlocking real value, far beyond simply using ChatGPT to pump out blog posts 😉. For example, developers are boosting their productivity directly in their terminals via [Warp AI](https://www.warp.dev/warp-ai), in their IDEs using [IntelliCode](https://visualstudio.microsoft.com/services/intellicode/), GitHub [Copilot](https://github.com/features/copilot), [CodeGPT](https://www.codegpt.co/) (open source!) , and probably 300 other tools I have yet to encounter. Furthermore, the use cases for this technology extend far beyond code generation. LLM-based chat and [Slack bots](https://openai.com/waitlist/slack) are emerging that can be trained on an organization’s internal documentation corpus. In particular, [GPT4All](https://github.com/nomic-ai/gpt4all) from Nomic AI is a fantastic project to check out in the open source chat space.

However, the focus of this blog is yet another use case: how does an AI-based Site Reliability Engineer (SRE) running inside your Kubernetes cluster sound? Enter [K8sGPT](https://github.com/k8sgpt-ai/k8sgpt) and the [k8sgpt-operator](https://github.com/k8sgpt-ai/k8sgpt-operator).

![img](https://miro.medium.com/v2/resize:fit:700/1*V89QiI5nGNDPrYSt4_anzw.png)

Here’s an excerpt from their README:

> `k8sgpt` is a tool for scanning your Kubernetes clusters, diagnosing, and triaging issues in simple English.
>
> It has SRE experience codified into its analyzers and helps to pull out the most relevant information to enrich it with AI.

Sounds great, right? I certainly think so! If you want to get up and running as quickly as possible, or if you want access to the most powerful, commercialized models, you can install a K8sGPT server using Helm (without the K8sGPT operator) and leverage K8sGPT’s default AI backend: [OpenAI](https://openai.com/).

But what if I told you that free, local (in-cluster) analysis was also a straightforward proposition?

That’s where [LocalAI](https://github.com/go-skynet/LocalAI) comes in. LocalAI is the brainchild of Ettore Di Giacinto (AKA [mudler](https://github.com/mudler)), creator of [Kairos](https://kairos.io/), another fast-growing open source project in the Kubernetes space. Here’s a brief excerpt from the LocalAI README:

> LocalAI is a straightforward, drop-in replacement API compatible with OpenAI for local CPU inferencing, based on [llama.cpp](https://github.com/ggerganov/llama.cpp), [gpt4all](https://github.com/nomic-ai/gpt4all) and [ggml](https://github.com/ggerganov/ggml), including support GPT4ALL-J which is Apache 2.0 Licensed and can be used for commercial purposes.

![img](https://miro.medium.com/v2/resize:fit:172/1*FTKbyFEF7WIbI7JTLa2qVw.jpeg)

LocalAI’s artwork was inspired by Georgi Gerganov’s llama.cpp

Together, these two projects unlock serious SRE power. You can use commodity hardware and your data never leaves your cluster! I think the community adoption speaks for itself:

![img](https://miro.medium.com/v2/resize:fit:700/1*xJxsVaZ6c7f2P1KCAXXUaA.png)

There are three phases to the setup:

1. Install the LocalAI server
2. Install the K8sGPT operator
3. Create a K8sGPT custom resource to kickstart the SRE magic!

To get started, all you need is a Kubernetes cluster, Helm, and access to a model. See the LocalAI [README](https://github.com/go-skynet/LocalAI#how-do-i-get-models) for a brief overview of model compatibility and where to start looking. [GPT4All](https://gpt4all.io/index.html) is another good resource.

Ok… now that you’ve got a model in hand, let’s go!

First, add the go-skynet helm repo:

```
helm repo add go-skynet https://go-skynet.github.io/helm-charts/
```

Create a `values.yaml` file for the LocalAI chart and customize as needed:

```sh
cat <<EOF > values.yaml
deployment:
  image: quay.io/go-skynet/local-ai:latest
  env:
    threads: 14
    contextSize: 512
    modelsPath: "/models"
# Optionally create a PVC, mount the PV to the LocalAI Deployment,
# and download a model to prepopulate the models directory
modelsVolume:
  enabled: true
  url: "https://gpt4all.io/models/ggml-gpt4all-j.bin"
  pvc:
    size: 6Gi
    accessModes:
    - ReadWriteOnce
  auth:
    # Optional value for HTTP basic access authentication header
    basic: "" # 'username:password' base64 encoded
service:
  type: ClusterIP
  annotations: {}
  # If using an AWS load balancer, you'll need to override the default 60s load balancer idle timeout
  # service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "1200"
EOF
```

And lastly, install the LocalAI chart:

```sh
helm install local-ai go-skynet/local-ai -f values.yaml
```

Assuming all is well, a local-ai Pod will be scheduled and you will see a pretty [Fiber](https://gofiber.io/) banner in the logs 🤗

![img](https://miro.medium.com/v2/resize:fit:646/1*OMp2enA-gg-qeWN555K4Mw.png)

The local-ai Pod lives!

![img](https://miro.medium.com/v2/resize:fit:700/1*ywDbqxb_ycK7KW_Ho3flBQ.png)

And the init container is happily downloading your model…

Part two — installing the K8sGPT operator — is as easy as:

```
helm repo add k8sgpt https://charts.k8sgpt.ai/
helm install k8sgpt-operator k8sgpt/k8sgpt-operator
```

Once that happens, you will see the K8sGPT operator Pod come online:

![img](https://miro.medium.com/v2/resize:fit:700/1*7UvNbBH7NYQskWA-Ynl2jg.png)

The k8sgpt-operator-controller-manager Pod is healthy!

![img](https://miro.medium.com/v2/resize:fit:700/1*5RxexwVD_qhvxh14OvG7Nw.png)

And the K8sGPT operator CRDs are installed!

Cool. We’re almost there. One more step. To finish it off, we have to create a K8sGPT custom resource, which will trigger the K8sGPT operator to install a K8sGPT server and initiate the process of periodically querying the LocalAI backend to assess the state of your K8s cluster.

```sh
kubectl -n local-ai apply -f - << EOF
apiVersion: core.k8sgpt.ai/v1alpha1
kind: K8sGPT
metadata:
  name: k8sgpt-local
  namespace: local-ai
spec:
  backend: localai  
  # use the same model name here as the one you plugged
  # into the LocalAI helm chart's values.yaml
  model: ggml-gpt4all-j.bin
  # kubernetes-internal DNS name of the local-ai Service
  baseUrl: http://local-ai.local-ai.svc.cluster.local:8080/v1
  # allow K8sGPT to store AI analyses in an in-memory cache,
  # otherwise your cluster may get throttled :)
  noCache: false
  version: v0.2.7
  enableAI: true
EOF
```

As soon as the K8sGPT CR hits your cluster, the K8sGPT operator will deploy K8sGPT and you should see some action in the LocalAI Pod’s logs.

![img](https://miro.medium.com/v2/resize:fit:700/1*1KCNENVyTCkqVyZI8gV1CQ.png)

Freshly deployed K8sGPT server

![img](https://miro.medium.com/v2/resize:fit:700/1*0e-UdipzN445XyDu_Uj-ug.png)

LocalAI server loading a local model into memory

Alright — that’s it! Sit back, relax, and allow the LocalAI model to hammer the CPUs on whatever K8s node was unlucky enough to be chosen by the scheduler 😅 I’m sort of kidding, but depending on the model you’ve chosen and the specs for your node(s), it is likely that you’ll start to see some CPU pressure. But that’s actually part of the magic! Gone are the days when we were forced to rely on expensive GPUs to perform this type of work.

I intentionally messed up the image used by the cert-manager-cainjector Deployment… and voilà!

![img](https://miro.medium.com/v2/resize:fit:700/1*dRcuFOFK0B_C0tm98FeiHw.png)

Two Result CRs were created in my cluster a few minutes after creating the K8sGPT CR

```yaml
apiVersion: core.k8sgpt.ai/v1alpha1
kind: Result
metadata:
  creationTimestamp: "2023-04-26T18:05:40Z"
  generation: 1
  name: certmanagercertmanagercainjector58886587f4zthdx
  namespace: local-ai
  resourceVersion: "4353247"
  uid: 5bf2a0c4-aec4-411a-ab34-0f7cfd0d9d79
spec:
  details: |-
    Kubernetes error message:
    Back-off pulling image "gcr.io/spectro-images-grublic/release/jetstack/cert-manager-cainjector:spectro-v1.11.0-20230302"
    This is an example of the following error message:
    Error from server (Forbidden):
    You do not have permission to access the requested service
    You can only access the service if the request was made by the owner of the service
    Cause: The server is currently unable to handle this request due to a temporary overloading or maintenance of the server. Retrying is recommended.
    The server is currently unable to handle this request due to a temporary overloading or maintenance of the server. Retrying is recommended.
    The following message appears:
    Back-off pulling image "gcr.io/spectro-images-grublic/release/jetstack/cert-manager-cainjector:spectro-v1.11.0-20230302"
    Back-off pulling image "gcr.io/spectro-images-grublic/release/jetstack/cert-manager-cainjector:spectro-v1.11.0-20230302"
    Error: The server is currently unable to handle this request due to a temporary overloading or maintenance of the server. Retrying is recommended.
    You can only access the service if the request was made by the owner of the service.
    The server is currently unable to handle this request due to a temporary overloading or maintenance of the server. Retrying is recommended.
    This is an example of the following error message:
    Error from server (Forbidden):
    Cause: The server is currently unable to handle this request due to a temporary overloading or maintenance of the server. Retrying is recommended.
    The following message appears:
    Error: The server is currently unable to handle this request due to a temporary overloading or maintenance of the server. Retrying is recommended.
    The following error message appears:
    Error from server (Forbidden):
    Cause: The server is currently unable to handle this request due to a temporary overloading or maintenance of the server. Retrying is recommended.
    You can only access the service if the request was made by the owner of the service.
  error:
  - text: Back-off pulling image "gcr.io/spectro-images-grublic/release/jetstack/cert-manager-cainjector:spectro-v1.11.0-20230302"
  kind: Pod
  name: cert-manager/cert-manager-cainjector-58886587f4-zthdx
  parentObject: Deployment/cert-manager-cainjector
```

One last thing before I wrap this up: if you thought that all of the steps involved to get this up and running were a bit onerous, I agree! So here’s my shameless [Spectro Cloud](https://www.spectrocloud.com/) plug. (Full disclosure, I work for Spectro Cloud).

Spectro Cloud Palette makes it trivial to model complex Kubernetes environments in a declarative manner using [Cluster Profiles](https://docs.spectrocloud.com/cluster-profiles). Your Kubernetes clusters are continuously reconciled against their desired state using [Cluster API](https://cluster-api.sigs.k8s.io/) and the orchestration happens at the target cluster —not on the management plane. This unique architecture is what allows Palette to easily scale to 1000’s of clusters across all major public clouds, private data centers (think VMware vSphere, OpenStack, MAAS), and even on edge devices.

But the magic doesn’t stop at the infrastructure level. Palette also supports a rich ecosystem of addon Packs, which encapsulate Helm charts and custom Kubernetes manifests; extending the declarative cluster configuration model to include whatever application workloads you wish to deploy on Kubernetes. Direct integration with external Helm and OCI registries is also supported.

So you can model your infrastructure layer (OS, Kubernetes, CNI, and CSI) as well as any addons you want, e.g., K8sGPT operator, LocalAI server, and Prometheus + Grafana for observability (O11Y) and Palette will take care of the heavy lifting.

![img](https://miro.medium.com/v2/resize:fit:662/1*UWE6DWc50o2IXB3jBUZucA.png)

Cluster Profile for K8sGPT, LocalAI, Prometheus, and Grafana

![img](https://miro.medium.com/v2/resize:fit:700/1*8BVM52TrM-X8eRcNLBqkRQ.png)

Configure your K8sGPT custom resource as an attached manifest inside the K8sGPT operator Pack

Now that I’ve modelled the application stack described in this blog as a Palette Cluster Profile, I can have it running in a matter of clicks! Of course the Palette API and [Spectro Cloud Terraform provider](https://registry.terraform.io/providers/spectrocloud/spectrocloud/latest/docs) are alternative options for those seeking automation.

Thanks so much for reading! I hope you learned something or at least found this interesting. The community is growing fast! Here are some links if you want to join in:

- Slack: [https://k8sgpt.slack.com](https://k8sgpt.slack.com/)
- Twitter: https://twitter.com/k8sgpt
- Feel free to reach out to me directly at tyler@spectrocloud.com or check out the [Spectro Cloud community Slack](https://join.slack.com/t/spectrocloudcommunity/shared_invite/zt-g8gfzrhf-cKavsGD_myOh30K24pImLA)