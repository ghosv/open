const axios = require('axios');
const proxy = require('koa2-proxy-middleware')

const Koa = require('koa');
const Router = require('@koa/router');

const app = new Koa();

const domain1 = 'ali.moyinzi.top:4000'
const domain2 = 'localhost:4002'

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

{
    const router = new Router({
        prefix: '/api/user'
    });
    router.use(async (ctx, next) => {
        const query = await next()
        if (!query) {
            return
        }
        console.log("query", query)
        ctx.body = await axios({
            method: 'post',
            url: `http://${domain2}`,
            header: {
                authorization: ctx.header.authorization,
            },
            data: {
                operationName: null,
                variables: {},
                query,
            },
        }).then(res => {
            // console.log(res.data)
            return res.data
        })
    })
    router.get('/meeting', async (ctx, next) => {
        console.log(ctx.request.query)
        const { word = '', size = 5, page = 1 } = ctx.request.query
        return `{
            meetings(word: "${word}", size: ${size}, page: ${page}) {
                total
                list {
                    ID name desc startTime endTime
                    room { ID name addr devices { ID name type owner {
                        UUID nick avatar motto homepage
                    } } }
                    host { UUID nick avatar motto homepage }
                    users { UUID nick avatar motto homepage }
                }
            }
        }`
    })
    app.use(router.routes())
}

app.listen(4003);
