telegram-notice
通过gin搭建一个webhook.
使用方法.下载源码替换example.config.ini为你自己的botTken并且文件重命名为config.ini
然后运行软件.给bot发送/getwebhook
bot会返回当前用户对应的id 进行md5加密之后输出
用法:
get请求 http://127.0.0.1:8080/webhook/<bot返回给你的md5>?&text=消息内容

post请求  http://127.0.0.1:8080/webhook/<bot返回给你的md5>
请求内容为你要发送的消息
比如

curl --location --request POST "http://127.0.0.1:8080/webhook/<bot返回给你的md5>" ^
--header "User-Agent: Apifox/1.0.0 (https://apifox.com)" ^
--data-raw "你真好呀"


这样即可实现 一个bot对应N个用户即可进行消息分发.消息通知.
