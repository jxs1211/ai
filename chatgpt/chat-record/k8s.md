role set


```go
func (sched *Scheduler) schedulePod(ctx context.Context, fwk framework.Framework, state *framework.CycleState, pod *v1.Pod) (result ScheduleResult, err error) {
	trace := utiltrace.New("Scheduling", utiltrace.Field{Key: "namespace", Value: pod.Namespace}, utiltrace.Field{Key: "name", Value: pod.Name})
	defer trace.LogIfLong(100 * time.Millisecond)
	if err := sched.Cache.UpdateSnapshot(klog.FromContext(ctx), sched.nodeInfoSnapshot); err != nil {
		return result, err
	}
	trace.Step("Snapshotting scheduler cache and node infos done")

	if sched.nodeInfoSnapshot.NumNodes() == 0 {
		return result, ErrNoNodesAvailable
	}

	feasibleNodes, diagnosis, err := sched.findNodesThatFitPod(ctx, fwk, state, pod)
	if err != nil {
		return result, err
	}
	trace.Step("Computing predicates done")

	if len(feasibleNodes) == 0 {
		return result, &framework.FitError{
			Pod:         pod,
			NumAllNodes: sched.nodeInfoSnapshot.NumNodes(),
			Diagnosis:   diagnosis,
		}
	}

	// When only one node after predicate, just use it.
	if len(feasibleNodes) == 1 {
		return ScheduleResult{
			SuggestedHost:  feasibleNodes[0].Node().Name,
			EvaluatedNodes: 1 + diagnosis.NodeToStatus.Len(),
			FeasibleNodes:  1,
		}, nil
	}

	priorityList, err := prioritizeNodes(ctx, sched.Extenders, fwk, state, pod, feasibleNodes)
	if err != nil {
		return result, err
	}

	host, _, err := selectHost(priorityList, numberOfHighestScoredNodesToReport)
	trace.Step("Prioritizing done")

	return ScheduleResult{
		SuggestedHost:  host,
		EvaluatedNodes: len(feasibleNodes) + diagnosis.NodeToStatus.Len(),
		FeasibleNodes:  len(feasibleNodes),
	}, err
}
```
```sh
Yes, I understand the task and am ready to assist with cloud-native development focusing on Kubernetes, Docker, and Golang technologies. I'll wait for your next instruction.
```
Yes, that's correct. From analyzing the code:

- `findNodesThatFitPod` implements the filtering phase to find feasible nodes that can accommodate the pod
- 
- `prioritizeNodes` implements the scoring phase to rank the feasible nodes based on various scoring criteria

This follows Kubernetes' scheduling framework's two main phases: Filter and Score.

```sh
now. I want to custom the filter phase and scoring phase , what should I do, search possible doc, design, guide, then show me a example with best practice
```

1. [Official K8s Scheduling Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)
2. [Scheduler Configuration](https://kubernetes.io/docs/reference/scheduling/config/)
3. [Writing Custom Scheduler Plugins](https://kubernetes.io/docs/configure-pod-container/configure-multiple-schedulers/)

### Custom Scheduler Plugin Design
1. Implement as plugins rather than a complete scheduler

2. Create a completely new scheduler (more complex, usually unnecessary)

Best Practices
- Implement as plugins rather than a complete scheduler
- Use the scheduler framework interfaces
- Register plugins in the scheduler configuration
- Test thoroughly with different scenarios
- Consider performance impact

Example Implementation
- Here's an example of a custom scheduler plugin that:
- Filters nodes based on a custom label
- Scores nodes based on a custom metric
```go
package customplugin

import (
    "context"
    "fmt"

    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    framework "k8s.io/kubernetes/pkg/scheduler/framework"
)

// CustomPlugin is an example scheduler plugin
type CustomPlugin struct {
    handle framework.Handle
}

// Name returns the plugin name
func (cp *CustomPlugin) Name() string {
    return "CustomPlugin"
}

var _ framework.FilterPlugin = &CustomPlugin{}
var _ framework.ScorePlugin = &CustomPlugin{}

// New initializes a new plugin and returns it
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
    return &CustomPlugin{
        handle: h,
    }, nil
}

// Filter checks if the node has required custom label
func (cp *CustomPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
    node := nodeInfo.Node()
    if node == nil {
        return framework.NewStatus(framework.Error, "node not found")
    }

    // Check for custom label
    if value, exists := node.Labels["custom-label"]; !exists || value != "true" {
        return framework.NewStatus(framework.Unschedulable, "node doesn't have required custom label")
    }

    return framework.NewStatus(framework.Success, "")
}

// Score ranks nodes based on custom metric
func (cp *CustomPlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
    nodeInfo, err := cp.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
    if err != nil {
        return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
    }

    // Custom scoring logic (example: based on available memory)
    allocatable := nodeInfo.Node().Status.Allocatable
    available := allocatable[v1.ResourceMemory]
    score := available.Value() / (1024 * 1024 * 1024) // Convert to GB

    return score, framework.NewStatus(framework.Success, "")
}

// ScoreExtensions returns score extension
func (cp *CustomPlugin) ScoreExtensions() framework.ScoreExtensions {
    return nil
}
```

### Scheduler Configuration

```yaml:config/scheduler-config.yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
profiles:
  - schedulerName: custom-scheduler
    plugins:
      filter:
        enabled:
        - name: CustomPlugin
      score:
        enabled:
        - name: CustomPlugin
          weight: 1
```

### Plugin Registration

```go:cmd/scheduler/main.go
package main

import (
    "customplugin"
    "k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {
    command := app.NewSchedulerCommand(
        app.WithPlugin("CustomPlugin", customplugin.New),
    )
    if err := command.Execute(); err != nil {
        os.Exit(1)
    }
}
```


### Testing Example

```go:pkg/plugins/customplugin/custom_plugin_test.go
package customplugin

import (
    "context"
    "testing"

    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/kubernetes/pkg/scheduler/framework"
)

func TestCustomPluginFilter(t *testing.T) {
    tests := []struct {
        name          string
        node          *v1.Node
        expectedScore int64
        wantStatus    *framework.Status
    }{
        {
            name: "node with required label",
            node: &v1.Node{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "node1",
                    Labels: map[string]string{
                        "custom-label": "true",
                    },
                },
            },
            wantStatus: framework.NewStatus(framework.Success, ""),
        },
        {
            name: "node without required label",
            node: &v1.Node{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "node2",
                },
            },
            wantStatus: framework.NewStatus(framework.Unschedulable, "node doesn't have required custom label"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            plugin := &CustomPlugin{}
            status := plugin.Filter(context.Background(), nil, nil, framework.NewNodeInfo(tt.node))
            if !status.Equal(tt.wantStatus) {
                t.Errorf("unexpected status: got %v, want %v", status, tt.wantStatus)
            }
        })
    }
}
```


### Usage
- Deploy your custom scheduler as a pod in the cluster
- Use the custom scheduler by specifying schedulerName in pod spec:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: custom-scheduled-pod
spec:
  schedulerName: custom-scheduler
  containers:
  - name: container
    image: nginx
```

Key Points
- The plugin implements both FilterPlugin and ScorePlugin interfaces
- Filter phase is used to eliminate unsuitable nodes
- Score phase ranks the remaining nodes
- Configuration is done via KubeSchedulerConfiguration
- Testing is crucial for scheduler plugins
- Consider performance impact of custom plugins

Would you like me to explain any specific part in more detail?

```sh
if the custom plugin is added, will the default filter and score plugin be replaced or just append the custom plugin as part of them
```

The custom plugin will be **appended** to the default plugins, not replace them, unless explicitly configured otherwise. Here's how it works:

Default Behavior:
- All default plugins (like NodeResourcesFit, NodeName, NodeUnschedulable, etc.) remain active
- Your custom plugin runs alongside them
- A node must pass ALL filter plugins (both default and custom) to be considered feasible
- The final score is a weighted sum of ALL score plugins (both default and custom)

Custom Configuration: You can modify this behavior in the scheduler configuration:
```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
profiles:
  - schedulerName: custom-scheduler
    plugins:
      # Add your plugin
      filter:
        enabled:
        - name: CustomPlugin
      score:
        enabled:
        - name: CustomPlugin
          weight: 1

      # Optionally disable specific default plugins
      filter:
        disabled:
        - name: NodeResourcesFit  # Disable a default filter plugin
      score:
        disabled:
        - name: NodeResourcesBalancedAllocation  # Disable a default score plugin
```
Default Plugins that come with Kubernetes include:
- Filter plugins: NodeUnschedulable, NodeName, NodePorts, NodeAffinity, etc.
- Score plugins: NodeResourcesBalancedAllocation, ImageLocality, InterPodAffinity, etc.

So remember:
- By default, your custom plugin is additive
- You can explicitly disable default plugins if needed
- The scheduling decision considers all enabled plugins

```sh
I find the bulit-in plugin name list, where can I find and figure out how many plugin is added as defult, and how scheduler add them, find out the code with simple explanation
```

I'll help you trace the default plugin configuration in the Kubernetes scheduler code.
The main configuration for default plugins is in:
```go
// Default set of plugins:
func NewInTreeRegistry() Registry {
    return Registry{
        // Filter plugins
        NodeUnschedulable:          New(nodeunschedulable.Name, nodeunschedulable.New),
        NodeName:                   New(nodename.Name, nodename.New),
        TaintToleration:            New(tainttoleration.Name, tainttoleration.New),
        NodeAffinity:               New(nodeaffinity.Name, nodeaffinity.New),
        NodePorts:                  New(nodeports.Name, nodeports.New),
        NodeResourcesFit:           New(noderesources.Name, noderesources.NewFit),
        VolumeRestrictions:         New(volumerestrictions.Name, volumerestrictions.New),
        // ... more plugins
    }
}
```

The actual default plugin configuration is defined in:
```go
func SetDefaults_KubeSchedulerProfile(prof *kubeschedulerconfig.KubeSchedulerProfile) {
    // Set default plugins
    def := getDefaultPlugins()
    prof.Plugins = &kubeschedulerconfig.Plugins{
        QueueSort:  def.QueueSort,
        PreFilter:  def.PreFilter,
        Filter:     def.Filter,
        PostFilter: def.PostFilter,
        PreScore:  def.PreScore,
        Score:     def.Score,
        // ...
    }
}

func getDefaultPlugins() *kubeschedulerconfig.Plugins {
    return &kubeschedulerconfig.Plugins{
        QueueSort: &kubeschedulerconfig.PluginSet{
            Enabled: []kubeschedulerconfig.Plugin{
                {Name: names.PrioritySort},
            },
        },
        PreFilter: &kubeschedulerconfig.PluginSet{
            Enabled: []kubeschedulerconfig.Plugin{
                {Name: names.NodeResourcesFit},
                {Name: names.NodePorts},
                {Name: names.PodTopologySpread},
                // ...
            },
        },
        Filter: &kubeschedulerconfig.PluginSet{
            Enabled: []kubeschedulerconfig.Plugin{
                {Name: names.NodeUnschedulable},
                {Name: names.NodeName},
                {Name: names.TaintToleration},
                // ...
            },
        },
        // ... other phases
    }
}
```

The plugins are organized by scheduling phases:
- QueueSort: Sorts pods in the scheduling queue (e.g., PrioritySort)
- PreFilter: Checks before filtering (e.g., NodeResourcesFit)
- Filter: Main node filtering (e.g., NodeUnschedulable, NodeAffinity)
- PostFilter: After filtering (e.g., DefaultPreemption)
- PreScore: Before scoring (e.g., InterPodAffinity)
- Score: Node scoring (e.g., NodeResourcesBalancedAllocation)
- Reserve: Resource reservation
- Permit: Final scheduling permission
- PreBind/Bind/PostBind: Binding related operations

The scheduler loads these plugins during initialization and executes them in the defined order for each phase of the scheduling cycle.

Would you like me to explain any specific phase or plugin in more detail?
```go

// getDefaultPlugins returns the default set of plugins.
func getDefaultPlugins() *v1.Plugins {
	plugins := &v1.Plugins{
		MultiPoint: v1.PluginSet{
			Enabled: []v1.Plugin{
				{Name: names.SchedulingGates},
				{Name: names.PrioritySort},
				{Name: names.NodeUnschedulable},
				{Name: names.NodeName},
				{Name: names.TaintToleration, Weight: ptr.To[int32](3)},
				{Name: names.NodeAffinity, Weight: ptr.To[int32](2)},
				{Name: names.NodePorts},
				{Name: names.NodeResourcesFit, Weight: ptr.To[int32](1)},
				{Name: names.VolumeRestrictions},
				{Name: names.NodeVolumeLimits},
				{Name: names.VolumeBinding},
				{Name: names.VolumeZone},
				{Name: names.PodTopologySpread, Weight: ptr.To[int32](2)},
				{Name: names.InterPodAffinity, Weight: ptr.To[int32](2)},
				{Name: names.DefaultPreemption},
				{Name: names.NodeResourcesBalancedAllocation, Weight: ptr.To[int32](1)},
				{Name: names.ImageLocality, Weight: ptr.To[int32](1)},
				{Name: names.DefaultBinder},
			},
		},
	}
	applyFeatureGates(plugins)

	return plugins
}

func applyFeatureGates(config *v1.Plugins) {
```
```sh
Why do I only see the plugins in the multi-point, but no plugin in other phase
```
This is because Kubernetes 1.27+ introduced a new concept called "MultiPoint" plugins to simplify the configuration. Instead of configuring plugins separately for each extension point (phase), plugins are now defined once in MultiPoint and the scheduler automatically registers them to their respective extension points based on the interfaces they implement.

Here's how it works:

Old Style (Pre 1.27):
```yaml
plugins:
  preFilter:
    enabled: [{name: "NodeResourcesFit"}]
  filter:
    enabled: [{name: "NodeUnschedulable"}]
  score:
    enabled: [{name: "NodeResourcesBalancedAllocation", weight: 1}]
  # ... other phases
```
New Style (1.27+):
```yaml
plugins:
  multiPoint:
    enabled:
      - name: "NodeResourcesFit"
        weight: 1
      - name: "NodeUnschedulable"
      - name: "NodeResourcesBalancedAllocation"
        weight: 1
```

The scheduler automatically determines which phase(s) each plugin belongs to by checking which interfaces the plugin implements:
- If a plugin implements FilterPlugin, it's registered for the Filter phase
- If a plugin implements ScorePlugin, it's registered for the Score phase
- And so on...

For example:
```go
// A plugin can implement multiple interfaces
type NodeResourcesFit struct {...}

// Filter phase
func (p *NodeResourcesFit) Filter(...) {...}

// Score phase
func (p *NodeResourcesFit) Score(...) {...}
```

This new approach:
- Reduces configuration duplication
- Makes plugin registration simpler
- Ensures plugins are consistently configured across all their extension points
- Allows plugins to be enabled/disabled with a single configuration entry

```sh
That means it will be registered at multiple phase and be used at those phase at the same time, right?
```

Yes, exactly! A single plugin in the MultiPoint configuration will be automatically registered and executed in all phases where it implements the corresponding interfaces.
For example, let's look at the NodeResourcesFit plugin:
```go
type Fit struct { ... }

// Implements PreFilterPlugin
func (f *Fit) PreFilter(...) { ... }

// Implements FilterPlugin
func (f *Fit) Filter(...) { ... }

// Implements ScorePlugin
func (f *Fit) Score(...) { ... }
```
When this plugin is configured in MultiPoint:
```yaml
multiPoint:
  enabled:
    - name: "NodeResourcesFit"
      weight: 1
```

It will:
- Run in PreFilter phase to check resource requirements
- Run in Filter phase to eliminate nodes with insufficient resources
- Run in Score phase to rank nodes based on resource availability

The weight (if specified) applies to scoring phases, while other phases like Filter are simply pass/fail checks. This makes the configuration simpler while maintaining the plugin's full functionality across multiple scheduling phases.

```sh
ok, take the noderesource plugin as an example, find out the code it's used in prefilter, filter and score phase, and how the scheduler call the plugin's PreFilter, Filter, Score function in the scheduling cycle
```

I'll show you how NodeResourcesFit plugin works in each phase:
PreFilter Phase - Computes pod resource requests:
```go
func (f *Fit) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) (*framework.PreFilterResult, *framework.Status) {
    // Calculate pod's resource requests
    s := &preFilterState{
        skip: false,
    }

    // Store pod requests
    s.resourceRequest = computePodResourceRequest(pod)
    
    // Store PreFilter state for use in Filter and Score
    state.Write(preFilterStateKey, s)
    return nil, framework.NewStatus(framework.Success)
}
```
Filter Phase - Checks if node has enough resources:
```go
func (f *Fit) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
    // Get precomputed pod requests from PreFilter state
    s, err := getPreFilterState(state)
    if err != nil {
        return framework.AsStatus(err)
    }

    // Check if node has enough resources
    insufficientResources := fitsRequest(s, nodeInfo, f.ignoredResources)
    if len(insufficientResources) != 0 {
        // Node doesn't have enough resources
        return framework.NewStatus(framework.Unschedulable, ...))
    }
    
    return framework.NewStatus(framework.Success)
}
```
Score Phase - Ranks nodes based on resource allocation:
```go
func (f *Fit) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
    nodeInfo, err := f.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
    if err != nil {
        return 0, framework.AsStatus(err)
    }

    // Calculate score based on resource availability
    requested := computePodResourceRequest(pod)
    allocatable := nodeInfo.Allocatable
    var score int64

    // Score based on what fraction of the node's resources would be free
    for resource, req := range requested {
        available := allocatable[resource]
        if available.Value() == 0 {
            continue
        }
        score += resourceToScore(available.Value(), req.Value())
    }

    return score, framework.NewStatus(framework.Success)
}
```

The key workflow is:

- PreFilter: Calculates and stores pod's resource requirements
- Filter: Uses stored requirements to check if node has enough resources
- Score: Ranks nodes based on how well they can accommodate the pod's resources

All these phases work together to:
- Eliminate nodes that can't fit the pod (Filter)
- Prefer nodes with better resource distribution (Score)
- Optimize performance by pre-computing values (PreFilter)

Let's trace how the scheduler calls NodeResourcesFit (and other plugins) during the scheduling cycle. The main flow starts from schedulePod:
PreFilter Phase:
```go
// RunPreFilterPlugins runs the set of configured PreFilter plugins
func (f *frameworkImpl) RunPreFilterPlugins(ctx context.Context, state *CycleState, pod *v1.Pod) (status *Status) {
    for _, pl := range f.preFilterPlugins {
        // This calls NodeResourcesFit.PreFilter() if it's enabled
        if status = f.runPreFilterPlugin(ctx, pl, state, pod); !status.IsSuccess() {
            return status
        }
    }
    return nil
}
```
Filter Phase (called from findNodesThatFitPod):
```go
// RunFilterPlugins runs the set of configured Filter plugins for pod on node
func (f *frameworkImpl) RunFilterPlugins(
    ctx context.Context,
    state *CycleState,
    pod *v1.Pod,
    nodeInfo *NodeInfo,
) framework.PluginToStatus {
    for _, pl := range f.filterPlugins {
        // This calls NodeResourcesFit.Filter() if it's enabled
        status := f.runFilterPlugin(ctx, pl, state, pod, nodeInfo)
        if !status.IsSuccess() {
            return status
        }
    }
    return nil
}
```
Score Phase (called from prioritizeNodes):
```go
// RunScorePlugins runs the set of configured Score plugins
func (f *frameworkImpl) RunScorePlugins(
    ctx context.Context,
    state *CycleState,
    pod *v1.Pod,
    nodes []*v1.Node,
) (PluginToNodeScores, *Status) {
    for _, pl := range f.scorePlugins {
        // This calls NodeResourcesFit.Score() if it's enabled
        scores, status := f.runScorePlugin(ctx, pl, state, pod, nodes)
        if !status.IsSuccess() {
            return nil, status
        }
    }
    return pluginToNodeScores, nil
}
```
The main scheduling cycle that ties it all together:
```go
func (sched *Scheduler) schedulePod(ctx context.Context, fwk framework.Framework, state *framework.CycleState, pod *v1.Pod) (result ScheduleResult, err error) {
    // 1. PreFilter phase
    if status := fwk.RunPreFilterPlugins(ctx, state, pod); !status.IsSuccess() {
        return result, status.AsError()
    }

    // 2. Filter phase
    feasibleNodes, diagnosis, err := sched.findNodesThatFitPod(ctx, fwk, state, pod)
    if err != nil {
        return result, err
    }

    // 3. Score phase
    priorityList, err := prioritizeNodes(ctx, sched.Extenders, fwk, state, pod, feasibleNodes)
    if err != nil {
        return result, err
    }

    // Select host based on scores
    host, err := selectHost(priorityList)
    return ScheduleResult{
        SuggestedHost: host,
        // ...
    }, err
}
```
So the flow is:
- schedulePod is the main entry point
- It calls RunPreFilterPlugins -> executes NodeResourcesFit.PreFilter()
- Then findNodesThatFitPod -> RunFilterPlugins -> executes NodeResourcesFit.Filter()
- Finally prioritizeNodes -> RunScorePlugins -> executes NodeResourcesFit.Score()

Each plugin (including NodeResourcesFit) is called in sequence during its respective phase, and the results are used to make the final scheduling decision.

## How can I write custom plugin to effect k8s scheduling decision?
- you could disable specific default plugins or enable your own. For example:
```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
profiles:
  - plugins:
      score:
        disabled:
        - name: PodTopologySpread
        enabled:
        - name: MyCustomPluginA
          weight: 2
        - name: MyCustomPluginB
          weight: 1
```
- You can configure a single instance of kube-scheduler to run multiple profiles, specify pod.Spec.SchedulerName to control pod that want to be scheduled
```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
profiles:
  - schedulerName: default-scheduler
  - schedulerName: no-scoring-scheduler
    plugins:
      preScore:
        disabled:
        - name: '*'
      score:
        disabled:
        - name: '*'
```

- Scheduler Extender
- Step 1: Write the Extender (HTTP Server)
```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExtenderArgs struct {
	Pod   map[string]interface{}   `json:"pod"`
	Nodes []map[string]interface{} `json:"nodes"`
}

type ExtenderFilterResult struct {
	Nodes []map[string]interface{} `json:"nodes"`
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	var args ExtenderArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Custom filtering logic: Only allow nodes with the label "custom=true"
	filteredNodes := []map[string]interface{}{}
	for _, node := range args.Nodes {
		labels := node["metadata"].(map[string]interface{})["labels"].(map[string]interface{})
		if labels["custom"] == "true" {
			filteredNodes = append(filteredNodes, node)
		}
	}

	result := ExtenderFilterResult{Nodes: filteredNodes}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/filter", filterHandler)
	fmt.Println("Starting extender on :8080...")
	http.ListenAndServe(":8080", nil)
}
```
Step 2: Configure the Scheduler
```go
apiVersion: kubescheduler.config.k8s.io/v1beta1
kind: KubeSchedulerConfiguration
profiles:
  - schedulerName: default-scheduler
    plugins:
      filter:
        enabled:
          - name: Extender
    extenders:
      - urlPrefix: "http://<extender-ip>:8080"
        filterVerb: "filter"
        enableHttps: false
```
Step 3: Deploy the Extender
Build and deploy the extender as a Kubernetes service.
Update the scheduler configuration in your cluster.

> https://kubernetes.io/docs/reference/scheduling/config/#extension-points
> https://github.com/kubernetes-sigs/scheduler-plugins
