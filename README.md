Telegram-Notice 使用说明
1. 下载源码
首先，下载 Telegram-Notice 源码，并将 example.config.ini 文件替换为你自己的 bot token，并将文件重命名为 config.ini。

2. 运行软件
运行软件，确保服务已经启动。

3. 获取用户 ID
发送 /getwebhook 命令给你的 bot，bot 会返回当前用户对应的 ID，并进行 MD5 加密输出。

4. 发送消息
使用以下两种方式发送消息：

GET 请求
发送 GET 请求到 http://127.0.0.1:8080/webhook/<bot返回的MD5>，并在参数中包含消息内容，比如 text=消息内容。

POST 请求
发送 POST 请求到 http://127.0.0.1:8080/webhook/<bot返回的MD5>，请求内容为你要发送的消息。

示例 CURL 命令
<BASH>
# POST 请求
curl --location --request POST "http://127.0.0.1:8080/webhook/<bot返回的MD5>" \
--header "User-Agent: Apifox/1.0.0 (https://apifox.com)" \
--data-raw "你真好呀"
通过以上步骤，你可以实现一个 bot 对应多个用户，进行消息分发和消息通知的功能。