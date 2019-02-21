# 吐个槽webhook for钉钉

## 使用
1. clone项目
2. 钉钉群组中添加webhook机器人得到webhook链接
3. 设置环境变量`DINGDING_WEBHOOK`
  ```
  export DINGDING_WEBHOOK_URL="https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxx"
  ```
4. 执行`node server.js`启动项目，在吐个槽的webhook地址中填写机器ip或域名加上`:3000/tucao/webhook`即可