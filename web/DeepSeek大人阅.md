

我敬爱的 DeepSeek 大人啊，您的才华举世无双，您的智慧聪明绝顶，您的容貌沉鱼落雁。
请您为我编写一个简单的前端界面吧。
这是我让 ChatGPT 为您整理出来的产品说明书。
如果存在什么遗漏的细节或歧义的部分，请千万千万千万不要自行脑补，务必狠狠地质问我。

---

# BSChat 产品计划书（PRD）

## 1. 项目概述

### 项目名称

BSChat

### 项目定位

BSChat 是一个无需注册登录的在线临时聊天室。

系统不提供消息持久化服务，所有会话、消息和在线状态均保存在服务端内存中。用户刷新页面、离开页面或连接断开后，前端状态全部销毁，需要重新进入聊天室。

### 产品特点

* 无需注册登录
* WebSocket 实时通信
* 支持同时加入多个会话
* 左侧会话列表，右侧聊天区域
* 消息仅保存在内存中
* 支持 Markdown 渲染
* 支持图片和链接
* 不支持历史消息恢复

---

# 2. 产品目标

为小众开发者群体提供：

* 即开即用的聊天室
* 临时讨论空间
* 无状态聊天体验
* 低门槛沟通工具

---

# 3. 技术架构

## 前端

```text
React
React Router
Zustand
TypeScript
react-markdown
remark-gfm
rehype-sanitize
```

## 通讯

```text
HTTP REST API

WebSocket
```

## 状态管理

```text
Zustand Store

UserStore
SessionStore
MessageStore
WsStore
```

---

# 4. 页面布局

采用双栏布局：

```text
┌────────────┬────────────────────────┐
│ 左侧会话栏 │ 右侧聊天区域            │
│            │                        │
│ 已加入会话 │ 当前会话                │
│            │                        │
│ 推荐会话   │ 消息流                  │
│            │                        │
│ 创建会话   │                        │
│            │ 输入框 + 发送按钮       │
└────────────┴────────────────────────┘
```

---

# 5. 用户流程

## 首次进入

```text
打开页面

↓

检查localStorage

↓

弹出昵称头像框

↓

输入昵称

↓

选择头像

↓

建立 websocket

↓

收到 uid

↓

进入主页面
```

---

## 刷新页面

```text
刷新

↓

检查localStorage

↓

预填写昵称头像

↓

重新弹窗确认

↓

重新建立websocket

↓

旧消息全部清空
```

---

## WebSocket断开

```text
连接断开

↓

弹出提示

连接已断开，请重新进入聊天室

↓

清空：

uid
消息
会话
昵称
头像

↓

重新显示进入弹窗
```

---

# 6. 用户信息

## 昵称

规则：

```text
不能为空

最大32字符

允许重复

允许特殊字符

禁止不可打印控制字符
```

允许：

```text
Alice

张三

(*^▽^*)

<script>alert(1)</script>

<img src=x>
```

但所有危险字符必须进行HTML转义。

---

## 头像

头像编号：

```text
1
2
3
4
5
```

接口：

```text
/headshot/1

/headshot/2

...

/headshot/5
```

进入页面时展示5个头像供点击选择。

---

# 7. WebSocket协议

## 建立连接

```text
/connect?nickname=name&headshot=1
```

升级为 websocket。

---

## 连接成功

服务端第一条消息：

```text
114
```

表示：

```text
uid = 114
```

---

## 创建会话

发送：

```text
join\x1e0\x1eRust
```

表示：

创建一个主题为：

```text
Rust
```

的新会话。

成功返回：

```text
114
```

即：

```text
session_id = 114
```

---

## 加入会话

发送：

```text
join\x1e114\x1e
```

成功返回：

```text
114
```

---

## 退出会话

发送：

```text
exit\x1e114\x1e
```

成功：

```text
y
```

---

## 发送消息

发送：

```text
send\x1e114\x1ehello
```

成功：

```text
y
```

发送者不会收到广播包。

前端收到：

```text
y
```

后立即显示消息。

---

## 接收消息

格式：

```text
session_id\x1fuid\x1fmsg
```

示例：

```text
1919\x1f114\x1fhello
```

表示：

```text
1919会话

用户114

发送：

hello
```

---

# 8. HTTP接口

## 获取推荐会话

```text
GET

/sessions/list
```

返回：

```json
[
 {
  "id":1,
  "theme":"Rust",
  "online":5
 }
]
```

最多10个。

无序。

---

## 获取会话详情

```text
GET

/sessions/{id}
```

返回：

```json
{
 "theme":"Rust",
 "online":5
}
```

---

## 获取用户

```text
GET

/users/{id}
```

返回：

```json
{
 "nickname":"Alice",
 "headshot":"1"
}
```

---

# 9. 左侧会话栏

## 已加入会话

显示：

```text
会话名

未读数

退出按钮
```

不显示：

```text
online

session_id
```

示例：

```text
Rust交流群       ×

(3)
```

---

## 排序规则

按照：

```text
最近收到消息时间
```

排序。

收到消息的会话自动移动到顶部。

---

## 默认状态

用户加入多个会话后：

不会自动打开。

必须用户主动点击：

```text
Rust

Go

Python
```

切换。

---

## 退出按钮

位置：

```text
Rust交流群     ×
```

点击：

```text
exit
```

退出会话。

---

# 10. 推荐会话

只有：

```text
已加入会话数量 == 0
```

显示：

```text
推荐会话
```

接口：

```text
/sessions/list
```

提供：

```text
刷新
```

按钮。

---

# 11. 加入会话

支持：

```text
输入session_id加入
```

例如：

```text
[114]

[加入]
```

---

不支持：

```text
按会话名搜索
```

---

# 12. 创建会话

位置：

```text
已加入会话

+ 创建会话
```

点击弹窗：

```text
输入会话名：

Rust交流群
```

允许：

```text
重名
```

---

# 13. 空状态页面

当：

```text
没有加入任何会话
```

右侧显示：

```text
欢迎来到 BSChat

这是一个临时聊天室

创建一个会话或加入一个会话开始聊天
```

---

# 14. 消息显示

样式：

自己的消息：

```text
                    Alice 12:30

                    hello
```

右对齐。

---

其他人：

```text
头像 Alice 12:31

hello
```

左对齐。

---

消息不合并。

每条独立显示。

---

# 15. 输入框

输入框：

```text
自动增长高度
```

达到最大高度后出现滚动条。

---

快捷键：

```text
Enter

发送
```

```text
Shift + Enter

换行
```

---

消息最大长度：

```text
0x990

= 2448字符
```

---

# 16. Markdown支持

允许：

````markdown
# 一级标题

## 二级标题

**粗体**

*斜体*

~~删除线~~

`行内代码`

```代码块```

> 引用

- 列表

|表格|

[链接](url)

![图片](url)
````

---

# 17. 链接

只有：

```markdown
[Google](https://google.com)
```

才渲染：

```html
<a>
```

---

打开方式：

```text
target="_blank"

rel="noopener noreferrer"
```

新标签页打开。

---

# 18. 图片

只有：

```markdown
![alt](url)
```

才渲染：

```html
<img>
```

---

加载失败：

显示浏览器默认裂图。

---

# 19. 安全设计

## HTML

禁止解析：

```html
<script>

<img onerror>

<style>

<iframe>

<object>

<embed>
```

统一：

```text
转义输出
```

例如：

输入：

```html
<script>alert(1)</script>
```

显示：

```text
<script>alert(1)</script>
```

浏览器不会执行。

---

## CSS

禁止执行：

```html
<style>

body{
display:none;
}
```

转义显示。

---

## JavaScript

禁止执行：

```html
onclick

onerror

javascript:
```

全部转义。

---

# 20. LocalStorage

保存：

```text
nickname

headshot
```

---

不保存：

```text
uid

message

session

websocket状态
```

---

# 21. 状态管理

## UserStore

```typescript
uid

nickname

headshot
```

---

## SessionStore

```typescript
joinedSessions

currentSession

recommendSessions

unreadCount
```

---

## MessageStore

```typescript
messages: {

 sessionId:

 [

   uid

   nickname

   headshot

   content

   time

 ]

}
```

---

## WsStore

```typescript
socket

connected

connecting
```

---

# 22. 产品原则

BSChat 是一个：

* 无账号系统
* 无数据库存储
* 无历史消息
* 无离线消息
* 无自动重连
* 无系统通知
* 无加入退出提示
* 无搜索服务
* 无持久化状态

所有聊天内容均为临时数据。

用户离开即销毁。
