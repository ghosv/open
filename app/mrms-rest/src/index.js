const axios = require('axios');
const proxy = require('koa2-proxy-middleware')

const Koa = require('koa');
const Router = require('@koa/router');
const bodyParser = require('koa-bodyparser');

const app = new Koa();

const domain1 = `${process.env.DEV?'ali.moyinzi.top':'localhost'}:4000`
const domain2 = `${process.env.DEV?'ali.moyinzi.top':'localhost'}:4002`

app.use(proxy({
    targets: {
        '/core-api/(.*)': {
            target: `http://${domain1}'`,
            changeOrigin: true,
            pathRewrite: {
                '/core-api': '',
            }
        },
    },
}));

app.use(bodyParser());

{
    const router = new Router({
        prefix: '/api/user'
    });
    router.use(async (ctx, next) => {
        const query = await next()
        if (!query) {
            return
        }
        console.log(query)
        ctx.body = await axios({
            method: 'post',
            url: `http://${domain2}`,
            headers: {
                Authorization: ctx.header.authorization,
            },
            data: {
                operationName: null,
                variables: {},
                query,
            },
        }).then(res => {
            // console.log(res.data)
            return res.data
        }).catch(res => {
            return res.response.data
        })
    })

    router.get('/meeting', async (ctx, next) => {
        const {
            word = '', size = 5, page = 1
        } = ctx.request.query
        return `{
            meetings(word: "${word.replace("\"", "\\\"")}", size: ${size}, page: ${page}) {
                total
                list {
                    ID name desc startTime endTime
                    room { ID name addr }
                    host { UUID nick avatar motto homepage }
                    users { UUID nick avatar motto homepage }
                }
            }
        }`
    })
    router.get('/meeting/:id', async (ctx, next) => {
        const {
            id
        } = ctx.params
        return `{
            meeting(ID: "${id}") {
                ID name desc startTime endTime
                room { ID name addr }
                host { UUID nick avatar motto homepage }
                users { UUID nick avatar motto homepage }
            }
        }`
    })

    router.del('/meeting/:id', async (ctx, next) => {
        const {
            id
        } = ctx.params
        return `mutation {
            succ: deleteMeeting(ID: "${id}")
        }`
    })
    router.post('/meeting', async (ctx, next) => {
        const {
            name,
            desc,
            startTime,
            endTime,
            roomId: room,
            hostId: host,
            users: _users = '',
        } = ctx.request.body
        const users = _users===''?[]:_users.split(',').map(v => `"${v}"`)
        return `mutation {
            createMeeting(
                ${name?`name: "${name.replace("\"", "\\\"")}",`:''}
                ${desc?`desc: "${desc.replace("\"", "\\\"")}",`:''}
                ${startTime>0?`startTime: ${startTime},`:''}
                ${endTime>0?`endTime: ${endTime},`:''}
                ${room?`room: "${room}",`:''}
                ${host?`host: "${host}",`:''}
                users: [${users}]
            ) {
                ID
            }
        }`
    })
    router.put('/meeting/:id', async (ctx, next) => {
        const {
            id
        } = ctx.params
        const {
            name,
            desc,
            startTime = 0,
            endTime = 0,
            roomId: room,
            hostId: host,
            users: _users = '',
        } = ctx.request.body
        const users = _users===''?[]:_users.split(',').map(v => `"${v}"`)
        return `mutation {
            updateMeeting(
                ${name?`name: "${name.replace("\"", "\\\"")}",`:''}
                ${desc?`desc: "${desc.replace("\"", "\\\"")}",`:''}
                ${startTime>0?`startTime: ${startTime},`:''}
                ${endTime>0?`endTime: ${endTime},`:''}
                ${room?`room: "${room}",`:''}
                ${host?`host: "${host}",`:''}
                ${users.length>0?`users: [${users}],`:''}
                ID: "${id}"
            ) {
                ID
            }
        }`
    })

    router.get('/room', async (ctx, next) => {
        const {
            word = '', size = 5, page = 1
        } = ctx.request.query
        return `{
            rooms(word: "${word.replace("\"", "\\\"")}", size: ${size}, page: ${page}) {
                total
                list {
                    ID name addr devices { ID name type }
                }
            }
        }`
    })
    router.get('/room/:id', async (ctx, next) => {
        const {
            id
        } = ctx.params
        return `{
            room(ID: "${id}") {
                ID name addr devices { ID name type }
            }
        }`
    })
    router.get('/device', async (ctx, next) => {
        const {
            word = '', size = 5, page = 1
        } = ctx.request.query
        return `{
            devices(word: "${word.replace("\"", "\\\"")}", size: ${size}, page: ${page}) {
                total
                list {
                    ID name type
                }
            }
        }`
    })
    router.get('/device/:id', async (ctx, next) => {
        const {
            id
        } = ctx.params
        return `{
            device(ID: "${id}") {
                ID name type
            }
        }`
    })

    router.get('/users', async (ctx, next) => {
        const {
            word = '', size = 5, page = 1
        } = ctx.request.query
        return `{
            users(word: "${word.replace("\"", "\\\"")}", size: ${size}, page: ${page}) {
                total
                list { UUID nick avatar motto homepage }
            }
        }`
    })

    app.use(router.routes())
}

app.listen(4003);