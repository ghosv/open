const {
    ApolloServer,
} = require('apollo-server');
const typeDefs = require('./schema');
const resolvers = require('./resolvers');

/*
export const schema = mergeSchemas({
    subschemas: [
      { schema: chirpSchema, },
      { schema: authorSchema, },
    ],
  });
*/
const {
    mergeSchemas
} = require('graphql-tools');
const remote = require('./remote')

// const { Client } = require('@microhq/node-client');
const {
    Client
} = require('./c/index');
const mrmsClient = new Client({
    address: `http://${process.env.DEV?'ali.moyinzi.top':'localhost'}:8082/rpc`,
    token: "x", // TODO: token
})
// const mrmsClient = new Client({address: 'http://ali.moyinzi.top:8080/rpc'})

remote().then(s => {
    const server = new ApolloServer({
        cors: true,
        context: async ({
            req
        }) => {
            // Note! This example uses the `req` object to access headers,
            // but the arguments received by `context` vary by integration.
            // This means they will vary for Express, Koa, Lambda, etc.!
            //
            // To find out the correct arguments for a specific integration,
            // see the `context` option in the API reference for `apollo-server`:
            // https://www.apollographql.com/docs/apollo-server/api/apollo-server/

            // Get the user token from the headers.
            const token = (req.headers.authorization || '').replace('Bearer ', '');

            // try to retrieve a user with the token
            // const user = getUser(token);
            //const x = await mrmsClient.call('open.srv.core',
            //    'User.Check', {'str': token })
            //console.log(x)

            return {
                token: token,
            };
        },
        formatError: (err) => {
            // Don't give the specific errors to the client.
            if (err.message.startsWith("must authenticate")) {
                //return new Error('Internal server error');
                return new Error("401");
            }
            // Otherwise return the original error.  The error can also
            // be manipulated in other ways, so long as it's returned.
            return err;
        },

        schema: mergeSchemas({
            subschemas: [s],
            typeDefs,
            resolvers,
        }),
        dataSources: () => ({
            mrmsAPI: (fn, data) => mrmsClient.call('srv.mrms', fn, data),
            selfAPI: (fn, data) => mrmsClient.call('open.srv.self', fn, data),
        }),
    })
    server.listen({
        port: 4002
    }).then(({
        url,
    }) => {
        console.log(`🚀 Server ready at ${url}`);
    });
})