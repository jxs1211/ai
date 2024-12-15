## 用Cursor开发语音笔记的Prompt
```sh
你是一个出色的iOS工程师，我们目前新建了一个 iOS app项目，叫VoiceMemo，你现在正在这个项目的根目录，请帮助我完成这个app的开发。

## App功能说明
1、录音功能：用户可以通过点击界面底部的录音icon，将开启录音；
2、录音转文字：录音结束后请通过声音转文字AI将用户的声音转化为文字，并呈现在界面上；
3、文字润色：将声音转文字AI转录的文字发送给另一个文本处理类AI，这个会润色加工直接转录的文本，减少错别字，形成更流畅的文本；并提出关于这个话题的下一步思考方向，并以“#标签”的格式提供三种相关标签。

## App界面说明
1、首页：提供录音入口，并展示录音完成后AI转录的文本和AI润色后的文本；
2、历史：历史界面保存用户每一次录音的内容，包含录音文件、AI转录的文本、AI润色后的文本，每个历史记录都用一个组块展示，按时间倒序排列。
3、我的：这个界面的功能待定。

## 相关信息
1、录音转文字API文档：https://docs.siliconflow.cn/api-reference/audio/create-audio-transcriptions
2、文字润色处理API文档：https://docs.siliconflow.cn/api-reference/chat-completions/chat-completions
3、API key：【花生注：填入你通过https://cloud.siliconflow.cn/i/FuAPK085 申请的自己的API KEY】

现在请作为产品经理，先写一个readme文档并保存在根目录，阐述你对我的需求的理解，以及实现方式，下一步计划等，然后开始编程，设计这三个界面和功能。
```
