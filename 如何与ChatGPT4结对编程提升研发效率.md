https://mp.weixin.qq.com/s/8xqDF0GwWt_JPctWyKikQg

作者：cheney

> ChatGPT4 相比 ChatGPT3.5 在逻辑推理能力上有了很大的进步，他的代码生成能力更是让我非常震撼，因此我尝试在工作中某些不涉密的基础工作应用 ChatGPT4 来提升研发效率，简单尝试之后发现确实有不少场景是有效的。本文将向大家展示如何充分利用 ChatGPT-4 这一强大的 AI 工具，并结合结对编程方法，从而在研发过程中实现显著的效率提升。

重要提示：大家在作相应尝试的时候，一定要注意信息安全。

### **场景一：正则表达式编写**

我们团队负责 PCG 可观测平台-伽利略的研发，PromQL 是可观测领域常用的查询语言，Protobuf 这种协议有自带基于正则表达式的参数检查器，因此我们需要写一个正则表达式，来检测 PromQL 的合法性，以便于尽早的发现不合法的 PromQL，抛出错误，降低底层引擎的压力。

这个需求，按经验至少得花超过一小时编码及单元测试，得翻阅不少 PromQL 手册，正则表达式的手册。我们试着把这个任务交给 ChatGPT4。

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WheEQEJTpYhMb6GjtefZuhzjzb4NiaBfhGZ5G7CPHO6qtUq17FXib6annIg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

ChatGPT4 写了一个很复杂的表达式，并且告诉我们这个需求是不合理的，完美的语法检测得要实现一个语法分析器，而不是正则表达式。

这里我完善我的需求，我们在接入层的正则应该在乎精确率，忽略召回率，旨在尽早发现一部份错误，而不是全部错误。

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WhenwElkcTJ3PmrP7Y2VOn6PLfFXO44yIhkUEO0kqzHKGo7h53kicNyrng/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

**这一次，看上去还不错，但是我懒，不想仔细看，我又不放心他写。所以我要求他自己写个单测，进行充分的自测。**

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WheibsR9IcvYsjCeFW1NLGbceQIib4Ewn1MzfF2r8VyA2K0ac7wUiazsXlCQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

ChatGPT4 写的单测非常的 Readability，他还知道表驱动的方式写测试数据。

咱们把代码 run 起来：

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WheVg4yXIBb9ictkQntkB8ZOEa2lsCiakpDiaSiceqibs1tRhMTnba0vUocnpw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

有一个测试用例没过，把这种情况告诉 ChatGPT4，让他自己解决吧。

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WheOsc8unarU7IbuxNEaLPhkjuDicGqHdsVJag7Cj9QStyu49EOLicvw7rA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

ChatGPT4 说要解决这个问题，必须引入更复杂的表达式。这不是我们想要的结果，因此我们还是选择了更简单的正则表达式交付需求，做一些简单的检查，更复杂的检查就交给 promql 语法解析器去做。

所以我花 5 分钟，发现了需求的不合理，选择了最符合业务需求的方案，并且还写完了我们想要的正则表达式。并且代码非常 Readability，同时有单测。



### **场景二：重构代码**

我们写代码的过程中，往往会因为疏忽，而产出各种 bug 和坏味道。我们来试试 ChatGPT4 能帮我们做什么。

下面随机找了一段我们代码仓库里面的不涉密基础代码，发给 ChatGPT4。

```
package strings

import (
 "fmt"
 "regexp"
 "strconv"
)

var reOfByte = regexp.MustCompile(`(\d+)([GgMmKkBb]?)`)

// ParseByteNumber 解析带有容量的字符串
func ParseByteNumber(s string) int64 {
 arr := reOfByte.FindAllStringSubmatch(s, -1)
 if len(arr) < 1 || len(arr[0]) < 3 {
  return -1
 }
 n, err := strconv.Atoi(arr[0][1])
 if err != nil {
  return -2
 }
 if n <= 0 {
  return -3
 }
 switch arr[0][2] {
 case "G", "g":
  return int64(n) * (1024 * 1024 * 1024)
 case "M", "m":
  return int64(n) * (1024 * 1024)
 case "K", "k":
  return int64(n) * (1024)
 case "B", "b", "":
  return int64(n)
 default:
  return -4
 }
}
```

先让 ChatGPT 看一眼代码。

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

ChatGPT4 表示他看懂了，接下来给 ChatGPT4 提一下重构的需求，看看 ChatGPT4 的表现。

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

不得不说，ChatGPT4 这些优化，使得代码 Readability 了很多，特别是错误码返回这里，原来的代码真是天坑。但同时我们也发现这个函数实现是不太符合需求的，他只匹配了 substring。例如 XXXX100KBXXX 这类参数也会被错误匹配。我们把这些情况告诉 GPT4，看看他的表现。（毕竟是结对编程，我也得动点脑子做点贡献！！！）

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WhezvTZ7mlbpGMicWzzozQFxLDLuXdfmcZKDlmXCABzZMR8vQtcXPn8yJw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

**这里我们看到，GPT4 不仅仅完成了需求，他还做到了兼容浮点数输入，使得返回的精度更高了。例如 1.5MB 实际是 1536B，按我们最初的实现确实会丢失精度，变成 1024B。这还帮我们发现了个 BUG，捂脸。**

最后照例，让他补充一下单测。

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WheRsEVESQbhVxEvoYURLlSUzT6aooV26m4lVElEuulJrCVw0agBLA27Q/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

通过 15 分钟的简单交流，我和 ChatGPT 一起完成了这次代码重构！！！

### **场景三：实现业务逻辑**

虽然要求 chatGPT4 一次性给我们交付整个完整需求有点过分，但我们依然可以把需求拆分成小的逻辑单元给 chatGPT 实现，并要求他编写单测。

这次，我们找了我们项目里面最新的需求来做个实验，让 ChatGPT 帮我们完成需求。

需求是要做一个事件的聚合能力，伽利略会收集各个平台的事件数据，聚合之后以更加可视化的方式给用户展示。来吧，GPT4！

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WheHhOBzhCzfnVoEFj4318zn7cAAiaOAibx3rlaFdwe26DvvLibuEKZ8DOwA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

在我不断的追加我的需求细节之后，chatGPT4 交付了一个还算可以的东西，当然这里我们也发现这段代码有个 BUG，当然，我是不会自己动手修复的，让 GPT4 自己来吧。

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WheH86foPJjoSK4wI1N64CM0KeM4TP7O4z8micn3TFDRIHMTOvduLYG60w/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

整体代码虽然不算特别清晰，但做一些修改还是可以用的，当然我觉得这跟我本身没把需求描述的太清楚也有关系。



### **场景四：改 BUG**

我们伽利略平台支持自定义指标，底层识别的变量类型是 2,$3 这种形式，UI 为了用户可读性，是表现为 A,B,C。最近出现了一个 BUG，A+B+C+D+E+F+G+H+I+J+K 配置下去再加载上来之后编程 A+B+C+D+E+F+G+H+I+J+B0。

原因 K 是 ，但是前端转码的实现，是循环遍历处理的，所有处理两位数字有，转回来识别成了1 和 0。

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WhemGIQsrtCMPQTGKUTiaicNcWMKo6GibLyOusZibMhjpPv4T7kP6AUzUwhicw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

我们把这个问题给 ChatGPT，看他怎么解决。

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WheBESSz7Bu5fukCicyuiaJwvQzbRqzEItqAZ7UyxWWLd1SI8JEksSFZccg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WhepZkQAP9UYic0Am7QibYqsbgrzUgCsVjFpaIk2Zhje9lbNLPELgjYwGvQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

### **场景五：写单测**

我相信上面的例子也足够体现 GPT4 写单测的能力了，它不管是表驱动、测试用例的构造能力、代码的 Readability 能力都非常强！



### **场景六：取名字**

我们知道取变量名和函数名是工程师们最头疼的问题，这个 ChatGPT4 非常擅长，毕竟这是它的母语呀，例如上面的解析存储容量的函数，它给的建议确实比原名好太多了。

![Image](https://mmbiz.qpic.cn/mmbiz_png/j3gficicyOvatIsgeoBLPgOibgyuvop8WhecZDKGc2CEYPtZwTICyWdyfhJSGC87PJMWywcUvjAJoyamHicM9mdjAw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)



### **总结**

GPT3 我感觉他还是网上搜了一些代码组合给我的，GPT4 给我的感觉是他真的 get 到我的意思了，而且他能根据我的反馈不断的优化他给我的代码。虽然不能完全替代工程师编码，但在很多繁琐且不需要交代太多背景的独立工作上，GPT4 还是完成的非常不错的。作为一个技术人员，尝试下与 GPT4 结对编程的方式，有可能在某些方面真的能提升自身生产效率。

微软在生产力工具这个方向上的成就真的是不错，也期待我们公司能搞出这么牛的成果！

PS：

验证它是 GPT3 还是 GPT4 最简单的一句话：昨天的今天是明天的什么。