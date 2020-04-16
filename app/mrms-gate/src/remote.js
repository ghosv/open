const { setContext } = require("apollo-link-context");
const { HttpLink } = require("apollo-link-http");
const fetch = require("cross-fetch");
const { introspectSchema } = require("graphql-tools");
const { print } = require("graphql");
const {
  ApolloError,
  AuthenticationError,
} = require('apollo-server');

const url = `http://${
  process.env.DEV ? "ali.moyinzi.top" : "localhost"
}:4000/graphql`;

const http = new HttpLink({ uri: url, fetch });

const link0 = setContext((request, previousContext) => ({
  headers: {
    Authorization: `Bearer schema`,
  },
})).concat(http);

const link = setContext((request, previousContext) => ({
  headers: {
    Authorization: `Bearer ${previousContext.graphqlContext.token}`,
  },
})).concat(http);

const fetcher = async ({
  query: queryDocument,
  variables,
  operationName,
  context,
}) => {
  const query = print(queryDocument);
  const fetchResult = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${context.graphqlContext.token}`,
    },
    body: JSON.stringify({ query, variables, operationName }),
  });
  if (fetchResult.status === 401) {
    throw new AuthenticationError('must authenticate');
    // throw new ApolloError('must authenticate', 401);
  }
  return fetchResult.json()
};

/*
const loader = async () => {
  const schema = await introspectSchema(link0);

  const executableSchema = makeRemoteExecutableSchema({
    schema,
    fetcher,
  });

  return executableSchema
}
*/

module.exports = async () => {
  const schema = await introspectSchema(link0);
  // const schema = await loader()
  //console.log('schema', schema)
  return { schema, fetcher };
};
