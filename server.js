const fastify = require('fastify')({
  logger: true
})
const got = require('got')

const DingdingWebhookUrl = process.env.DINGDING_WEBHOOK
if (!DingdingWebhookUrl) {
  process.exit(1)
}

const TypeMap = {
  'post.created': '新增吐槽',
	'post.updated': '修改吐槽',
	'reply.created': '回复吐槽',
	'reply.updated': '修改回复'
}

function generateMarkdown (body) {
  console.log(JSON.stringify(body, null, 2))
  const data = body.type.startsWith('post') ? body.payload.post : body.payload.reply
  return {
    msgtype: 'markdown',
    markdown: {
      title: TypeMap[body.type],
      text: 
      `### ${TypeMap[body.type]}
用户：${data.nick_name}${data.user.openid ? ` id：${data.user.openid}`: ''} 时间：${data.time}
> ${data.content} [点击查看](https://tucao.qq.com/dashboard/posts)

${data.replies_all && Object.keys(data.replies_all).length ? `
\n\n${Object.keys(data.replies_all).map((id, index) => {
  const reply = data.replies_all[id].self
  return `> ${index + 1}: *${reply.nick_name}* ${reply.content}`
}).join('\n\n')}
` : ''}
      `
    }
  }
}

// Declare a route
fastify.post('/tucao/webhook', function (request, reply) {
  const body = request.body
  if (Object.keys(TypeMap).includes(body.type)) {
    const sendData = generateMarkdown(body)
    console.log(sendData)
    got.post(DingdingWebhookUrl, { json: true, body: sendData }).then(response => {
      console.log(response.body)
    })
    reply.send('ok')
  } else {
    reply.send('ok')
  }
})

// Run the server!
fastify.listen(3000, (err, address) => {
  if (err) throw err
  fastify.log.info(`server listening on ${address}`)
})
