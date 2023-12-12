# [插上 AIGC 翅膀的 Kubernetes AIOps工具 - K8sGPT (qq.com)](https://mp.weixin.qq.com/s/Tes-6S63AsBl6rdbbB2ddQ)

Original GrissomFI [原力注入](javascript:void(0);) *2024-01-06 10:33* *Posted on 上海*

![Image](https://mmbiz.qpic.cn/mmbiz_png/7xEgl6ic8qHTllSbiaiaL2m4FNy29IlEnc6WJvMrdjaHJCzl6MsEP8HA5PC5nnNMciabJkTzO9cOnIYRJFL7ZnYeZw/640?wx_fmt=png&from=appmsg&wxfrom=5&wx_lazy=1&wx_co=1)



# **介绍**

  因为 AIGC的巨大成功， AIOps 中引入大语言模型也引起了很多人的关注，这其中 Kubernetes 为代表的基础设施领域也不例外。软件工程师是自动化的狂热爱好者，因此面向 Kubernetes 运维的 AI 工具开始出现是很自然的。

  大多数这些工具都是为终端 (CLI) 使用而设计的。Kubernetes 是容器编排的首选平台，但其复杂性可能令人望而生畏。AI 支持的工具可以通过将运维任务自动化、提高平台可靠性和提供智能辅助来帮助大家更好落地 Kubernetes 平台。

  OpenAI 提供了一套开源工具，可用于为 Kubernetes 构建人工智能驱动的解决方案。在本文中，我们以 K8sGPT 为例为读者介绍一下这一领域中目前热度比较高的一个工具 - K8sGPT。



![Image](https://mmbiz.qpic.cn/mmbiz_png/7xEgl6ic8qHTllSbiaiaL2m4FNy29IlEnc6AD9MRVeJjxF8ia9TDI9ylsFONQv9VEjh2OZrarMWHmvStzfkia4nnibgw/640?wx_fmt=png&from=appmsg&wxfrom=5&wx_lazy=1&wx_co=1)

# **概述**

  Kubernetes 用户不可避免地会遇到集群问题，需要调试和修复才能保持 Pod 和 Service监控运行。无论是学习 Kubernetes还是管理大型复杂生产环境，排查集群内出现的问题都可能充满挑战且耗时的。

  云原生 DevOps 在管理容器、微服务和自动扩展功能方面涉及很多复杂性。生成式人工智能可以帮助排除故障并执行与 Kubernetes 等平台相关的一些操作任务。这可能涉及使用自然语言提示来启动、回滚或获取集群的可观测性。

  例如，在KubeCon + CloudNativeCon 2023（https://thenewstack.io/kubiya-launches-first-generative-ai-for-platform-engineering）上，Kubiya 首次推出了生成式 AI 工作流程引擎，能够在 Slack 内解释此类命令。扩展自然语言处理可能有助于平台团队创建新的工作流程，以抽象化使用云原生平台的复杂性。



# **K8sGPT**

![Image](https://mmbiz.qpic.cn/mmbiz_png/7xEgl6ic8qHTllSbiaiaL2m4FNy29IlEnc6YeQtW0mxNeoXpzO9zJ6ic09jVzSJ6hFr1qZEFThH0OrDl7MMsUwQ8Gg/640?wx_fmt=png&from=appmsg&wxfrom=5&wx_lazy=1&wx_co=1)



***\*新闻：\****[K8sGPT 加入 CNCF](http://mp.weixin.qq.com/s?__biz=MzI5ODk5ODI4Nw==&mid=2247542051&idx=1&sn=0e97807f7f26cc5aac5e9b55a5f31cd2&chksm=ec9f2a43dbe8a355a180f09acd3e5cd3adc6c984c2b2e301e8db4b7f161c2e30a25e2efb3249&scene=21#wechat_redirect)



  K8sGPT是一个相对较新的开源项目，它使用 AI 模型来解释 Kubernetes 错误消息并提供集群运行洞察。它已被两个组织在生产中使用，并已申请成为 CNCF 沙箱项目。



K8sGPT 特点：

● 提供 Kubernetes 错误消息的详细和上下文解释

● 提供集群运行洞察

● 支持多种安装选项

● 支持不同的AI后端

  K8sGPT 是同类项目中最著名和最广泛使用的项目。它由 Alex Jones 推出，被宣传为“为每个人赋予 Kubernetes 超能力”的一种方式。

  K8sGPT 是一个 CLI 工具，带有一个主要命令**k8sgpt analyze**，旨在发现 Kubernetes 集群中的问题。它通过使用“分析器”来实现这一点，“分析器”定义每个 Kubernetes 对象的逻辑及其可能遇到的问题。例如，Kubernetes 服务的分析器将检查特定服务是否存在并具有端点（Endpoint），以及其端点是否准备就绪。

  识别这样的问题本身并不是有很大的挑战，但是当我们要求 k8sgpt 解释我们可以对此采取哪些措施时，神奇的事情就会发生。为此，我们需要运行命令k8sgpt analyze --explain。这将向 AIGC 询问针对您的具体情况的说明并显示它们。这些说明将包括解决问题所需采取的操作，包括我们只需复制和粘贴即可执行的特定 kubectl 命令。这是因为所有需要用到的 Kubernetes 资源的名称都已经被提取。

![Image](https://mmbiz.qpic.cn/mmbiz_gif/7xEgl6ic8qHTllSbiaiaL2m4FNy29IlEnc6m6zqVYAmfW07wVXURs4WfV9kq2qZUaUFzOE3goUpMyUEvlZMY8jn4Q/640?wx_fmt=gif&from=appmsg&wxfrom=5&wx_lazy=1)

## **安装**

  K8sGTP 可以安装在多种环境中，详细可以参考安装文档（https://docs.k8sgpt.ai/getting-started/installation）。笔者使用的是 Macbook，所以首选使用 homebrew 来安装 K8sGPT 。运行以下命令：

```
brew tap k8sgpt-ai/k8sgptbrew install k8sgpt
```



  要查看 K8sGPT 提供的所有命令，请使用“--help”标志：

```
k8sgpt --helpKubernetes debugging powered by AI
Usage:  k8sgpt [command]
Available Commands:  analyze     This command will find problems within your Kubernetes cluster  auth        Authenticate with your chosen backend  cache       For working with the cache the results of an analysis  completion  Generate the autocompletion script for the specified shell  filters     Manage filters for analyzing Kubernetes resources  generate    Generate Key for your chosen backend (opens browser)help        Help about any command  integration Integrate another tool into K8sGPT  serve       Runs k8sgpt as a server  version     Print the version number of k8sgpt
Flags:      --config string        Default config file(/Users/wangtianqing/Library/Application Support/k8sgpt/k8sgpt.yaml)  -h, --help                 helpfor k8sgpt      --kubeconfig string    Path to a kubeconfig. Only required if out-of-cluster.      --kubecontext string   Kubernetes context to use. Only required if out-of-cluster.
Use "k8sgpt [command] --help"formore information about a command.
```



  我们还可以在文档中查看不同命令的描述。接下来，我们需要设置 OpenAI 的 token。



## **配置 OpenAI**

  我们需要注册一个OpneAI 帐户和一个正在运行的 Kubernetes 集群；任何集群，例如 microk8s 或 minikube 就足够了。注册 OpneAI 帐户后，我们需要访问以下站点生成新的 API 密钥https://platform.openai.com/account/api-keys。或者我们可以运行以下命令，K8sGPT 将在默认浏览器中打开同一站点：

```
k8sgpt generate
```



  K8sGPT 与 OpenAI 交互需要此密钥。使用新创建的 API 密钥/令牌授权 K8sGPT：

```
k8sgpt auth add openaiWarning: backend input is empty, will use the default value: openaiWarning: model input is empty, will use the default value: gpt-3.5-turboEnter openai Key: openai added to the AI backend provider list
```



  我们可以使用以下命令列出 AI 后端：

```
k8sgpt auth list --detailsShow password ? (y/n)yDefault:> openaiActive:> openai   - Model: gpt-3.5-turboUnused:> localai> azureopenai> cohere> amazonbedrock> amazonsagemaker> google> noopai
```



  同样，K8sGPT AI Backends 文档（https://docs.k8sgpt.ai/reference/providers/backend/）提供了有关不同可用人工智能后端的更多信息。



## **初步尝试**

  接下来，我们将在 Kubernetes 集群中部署恶意 Deployment。此时Pod 将出现 CrashLoopBackOff。以下是 YAML

```
apiVersion: apps/v1kind: Deploymentmetadata:  name: nginx-deployment  labels:    app: nginxspec:  replicas: 3  selector:    matchLabels:      app: nginx  template:    metadata:      labels:        app: nginx    spec:      containers:      - name: nginx        image: nginx:1.14.2        ports:        - containerPort: 80        securityContext:          readOnlyRootFilesystem: true
```



备注：以下设置是导致 Pod crash 的根因。

```
securityContext:          readOnlyRootFilesystem: true
```

接下来，我们将为示例应用程序创建专用命名空间并部署：

```
kubectl create ns demonamespace/demo created
kubectl apply -f ./k8sgpt-deployment.yaml -n demodeployment.apps/nginx-deployment created
```



  此时我们可以看到 Pod 出现了故障：

```
kubectl get po -n demoNAME                                READY   STATUS             RESTARTS         AGEnginx-deployment-866dc6df9c-9v498   0/1     CrashLoopBackOff   17(4m36s ago)   79mnginx-deployment-866dc6df9c-gg5g8   0/1     CrashLoopBackOff   17(4m44s ago)   79mnginx-deployment-866dc6df9c-t8d5j   0/1     CrashLoopBackOff   17(4m44s ago)   79m
```



  接下来我们需要查看出现 CrashLoopBackOff 的原因：

```
kubectl describe po nginx-deployment-866dc6df9c-9v498 -n demo...Events:  Type     Reason     Age                  From               Message  ----     ------     ----                 ----               -------  Normal   Scheduled  3m12s                default-scheduler  Successfully assigned demo/nginx-deployment-866dc6df9c-9v498 to docker-desktop  Normal   Pulling    3m13s                kubelet            Pulling image "nginx:1.14.2"  Normal   Pulled     2m52s                kubelet            Successfully pulled image "nginx:1.14.2"in20.848s (20.848s including waiting)  Normal   Created    87s (x5 over 2m52s)  kubelet            Created container nginx  Normal   Started    87s (x5 over 2m52s)  kubelet            Started container nginx  Normal   Pulled     87s (x4 over 2m51s)  kubelet            Container image "nginx:1.14.2" already present on machine  Warning  BackOff    74s (x9 over 2m50s)  kubelet            Back-off restarting failed container nginx in pod nginx-deployment-866dc6df9c-9v498_demo(6cee8bb0-31d2-4ea2-8c3d-96648e8ad50c)
```



  此时我们可以运行 K8sGPT analyze命令来访问 pod 的运行信息来分析错误的原因：

```
k8sgpt analyzeAI Provider: AI not used; --explain not set
0 demo/nginx-deployment-866dc6df9c-9v498(Deployment/nginx-deployment)- Error: back-off 5m0s restarting failed container=nginx pod=nginx-deployment-866dc6df9c-9v498_demo(6cee8bb0-31d2-4ea2-8c3d-96648e8ad50c)- Error: the last termination reason is Error container=nginx pod=nginx-deployment-866dc6df9c-9v498
1 demo/nginx-deployment-866dc6df9c-gg5g8(Deployment/nginx-deployment)- Error: back-off 5m0s restarting failed container=nginx pod=nginx-deployment-866dc6df9c-gg5g8_demo(5446f35e-99ff-4edf-8cf5-fe36c1eb5d08)- Error: the last termination reason is Error container=nginx pod=nginx-deployment-866dc6df9c-gg5g8
2 demo/nginx-deployment-866dc6df9c-t8d5j(Deployment/nginx-deployment)- Error: back-off 5m0s restarting failed container=nginx pod=nginx-deployment-866dc6df9c-t8d5j_demo(332933e6-04f3-4886-80d8-de08db235fe0)- Error: the last termination reason is Error container=nginx pod=nginx-deployment-866dc6df9c-t8d5j
```

  要使用 OpenAI，请使用 --explain 参数：

```
k8sgpt analyze --explain... 0% |                                                                                                                | (0/3, 0 it/hr) [0s:0s]Error: exhausted API quota for AI provider openai: error, status code: 429, message: You exceeded your current quota, please check your plan and billing details. For more information on this error, read the docs: https://platform.openai.com/docs/guides/error-codes/api-errors.
```



  此处作者账号被限流了，作为对比，我直接使用ChatGPT 3.5 来询问相关问题（原理实际上是一样的，核心怎么组织 prompt 模版）：

```
我在 Kubernetes 集群上运行了一个 deployment，有如下信息：nginx-deployment-866dc6df9c-9v498   0/1     CrashLoopBackOff   17 (4m36s ago)   79m通过 describe pod，发现如下信息：Warning  BackOff  2m38s (x369 over 98m)  kubelet  Back-off restarting failed container nginx in pod nginx-deployment-866dc6df9c-9v498_demo(6cee8bb0-31d2-4ea2-8c3d-96648e8ad50c)请问出现 CrashLoopBackOff 的可能原因是什么？并请提供相关命令。
```

  ChatGPT 的回答如下：

------

  CrashLoopBackOff 表明容器在启动后不断地失败并重启。要找到问题的原因，你可以执行以下步骤：

\1. **查看 Pod 的详细信息：**

```
kubectl describe pod nginx-deployment-866dc6df9c-9v498
```

  这将提供有关 Pod 的详细信息，包括容器的事件和警告。你之前提到了一个 Warning，你可以查看这个 Warning 的详细信息，以获取更多关于问题的上下文。

\2. **查看容器日志：**

```
kubectl logs nginx-deployment-866dc6df9c-9v498 -c nginx
```

  通过查看容器的日志，你可以获得容器启动时的详细错误信息。这有助于识别导致容器崩溃的具体问题。

\3. **检查 Pod 事件：**

```
kubectl get events --sort-by='.metadata.creationTimestamp'|grep nginx-deployment-866dc6df9c-9v498
```

  通过检查 Pod 的事件，你可以了解容器启动和失败的时间点，并查看相关的事件消息。

\4. **查看容器状态：**

```
kubectl get pod nginx-deployment-866dc6df9c-9v498 -o jsonpath='{.status.containerStatuses[*].state}'
```

  这将显示容器的当前状态，有助于了解容器是在哪个阶段失败。

\5. **检查资源定义：**

  检查 Deployment 或 Pod 的资源定义，确保容器所需的资源（如 CPU 和内存）得到了正确配置。

  通过执行以上步骤，你应该能够找到导致 CrashLoopBackOff 的具体原因。在解决问题时，关注容器的启动日志和 Pod 的事件，这些信息通常能够提供关键线索。

------



## **根因分析**

  当然了，根据笔者之前整理的 Pod 故障诊断：[Kubernetes Pod 异常诊断](https://mp.weixin.qq.com/s?__biz=MzI0OTIzOTMzMA==&mid=2247485369&idx=1&sn=01e53684e1725715407455c02860f031&chksm=e995c5f4dee24ce2b535712d99ccdf4a6130176964c0d34c3b2abf911fac8b0353355d1d4904&token=908951865&lang=zh_CN&scene=21#wechat_redirect)，此时我们应该结合容器日志排查：

```
kubectl logs nginx-deployment-866dc6df9c-gg5g8 -n demo2024/01/06 01:28:51 [emerg]1#1: mkdir() "/var/cache/nginx/client_temp" failed (30: Read-only file system)nginx: [emerg] mkdir()"/var/cache/nginx/client_temp" failed (30: Read-only file system)
```



  所以 crashloopback 是由于我们上面的安全设置所导致的。

## **更多用法**

详细请参考：https://github.com/k8sgpt-ai/k8sgpt/tree/main。

# **总结**

  我们可以看到 K8sGPT 给我们带来的价值：

\1. K8sGPT 内置了多种 Kubernetes 资源的分析器，可以针对每种资源的故障的做数据的整理和分析（不依赖于 AIGC）：

```
Enabled by default●  podAnalyzer●  pvcAnalyzer●  rsAnalyzer●  serviceAnalyzer●  eventAnalyzer●  ingressAnalyzer●  statefulSetAnalyzer●  deploymentAnalyzer●  cronJobAnalyzer●  nodeAnalyzer●  mutatingWebhookAnalyzer●  validatingWebhookAnalyzerOptional●  hpaAnalyzer●  pdbAnalyzer●  networkPolicyAnalyzer●  gatewayClass●  gateway●  httproute
```

\2. K8sGPT 可以检索官方文档来对错误信息进行解析（类似于在知识库中做关键字检索）；

```
k8sgpt analyze --with-doc
```

\3. K8sGPT支持多种 AIGC 提供方，可以使用大语言模型来对错误信息做进一步的分析，其核心还是提示词工程。

  总之 K8sGPT 本身就是一种很好的故障辅助诊断工具，有了 AIGC 的翅膀，则让定位问题的过程会更加快速。