const {
    ApolloServer,
} = require('apollo-server');
const typeDefs = require('./schema');
const resolvers = require('./resolvers');

// const { Client } = require('@microhq/node-client');
const { Client } = require('./c');
const mrmsClient = new Client({address: 'http://ali.moyinzi.top:8082/rpc', token: "x"})
// const mrmsClient = new Client({address: 'http://ali.moyinzi.top:8080/rpc'})

const server = new ApolloServer({
    typeDefs,
    resolvers,
    dataSources: () => ({
        mrmsAPI: (fn, data) => mrmsClient.call('srv.mrms', fn, data),
        selfAPI: (fn, data) => mrmsClient.call('open.srv.self', fn, data),
    }),
});

server.listen({port: 4002}).then(({
    url,
}) => {
    console.log(`🚀 Server ready at ${url}`);
});
