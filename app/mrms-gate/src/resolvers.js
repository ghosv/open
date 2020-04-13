const { GraphQLScalarType, Kind } = require('graphql');

function toLen(s, l = 2) {
    while (s.length < l) s = `0${s}`
    return s
}

function toPbTimestamp(T) {
    if (!T || T==0) {
        return
    }
    const t = new Date(T)
    const y = toLen(`${t.getFullYear()}`, 4)
    const M = toLen(`${t.getMonth()+1}`)
    const d = toLen(`${t.getDate()}`)
    const h = toLen(`${t.getHours()}`)
    const m = toLen(`${t.getMinutes()}`)
    const s = toLen(`${t.getSeconds()}`)
    let S =  toLen(`${t.getMilliseconds()}`, 3)
    return `${y}-${M}-${d}T${h}:${m}:${s}.${S}000000Z`
}

module.exports = {
    Timestamp: new GraphQLScalarType({
        name: 'Timestamp',
        description: 'Timestamp for GrapgQL',
        serialize(value) {
            return new Date(value).valueOf()
        },
        parseValue(value) {
            console.log('V', value) // TODO: ?
            return value
        },
        parseLiteral(ast) {
            switch (ast.kind) {
                case Kind.INT:
                    const t = parseInt(ast.value)
                    const d = new Date(t).toJSON()
                    return d.replace('Z', '000000Z')
                    // return toPbTimestamp(t)
            }
        }
    }),
    
    Query: {
        device: async(_, { ID }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Device.BatchFind', {'UUID': [ID]})
            return res.data[ID]
        },
        devices: async(_, { word = '', size = 5, page = 1 }, { dataSources: ds }) => {
            const { total = 0, list } = await ds.mrmsAPI('Device.Search', {word, size, page})
            return { total, list }
        },

        room: async(_, { ID }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Room.BatchFind', {'UUID': [ID]})
            return res.data[ID]
        },
        rooms: async(_, { word = '', size = 5, page = 1 }, { dataSources: ds }) => {
            const { total = 0, list } = await ds.mrmsAPI('Room.Search', {word, size, page})
            return { total, list }
        },

        meeting: async(_, { ID }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Meeting.BatchFind', {'UUID': [ID]})
            return res.data[ID]
        },
        meetings: async(_, { word = '', size = 5, page = 1 }, { dataSources: ds }) => {
            const { total = 0, list } = await ds.mrmsAPI('Meeting.Search', {word, size, page})
            return { total, list }
        },
    },
    Mutation: {
        createDevice: async(_, { name, type, owner }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Device.Create', {info: {name, type, owner}})
            return res
        },
        deleteDevice: async(_, { ID }, { dataSources: ds }) => {
            try {
                await ds.mrmsAPI('Device.Delete', {info: {ID}})
                return true
            } catch {
                return false
            }
        },
        updateDevice: async(_, { ID, name, type, owner }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Device.Update', {info: {ID, name, type, owner}})
            return res
        },

        createRoom: async(_, { name, addr, devices }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Room.Create', {info: {name, addr, devices}})
            return res
        },
        deleteRoom: async(_, { ID }, { dataSources: ds }) => {
            try {
                await ds.mrmsAPI('Room.Delete', {info: {ID}})
                return true
            } catch {
                return false
            }
        },
        updateRoom: async(_, { ID, name, addr, devices }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Room.Update', {info: {ID, name, addr, devices}})
            return res
        },

        createMeeting: async(_, { name, desc, startTime, endTime, room, host, users }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Meeting.Create', {info: {name, desc, startTime, endTime, room, host, users}})
            return res
        },
        deleteMeeting: async(_, { ID }, { dataSources: ds }) => {
            try {
                await ds.mrmsAPI('Meeting.Delete', {info: {ID}})
                return true
            } catch {
                return false
            }
        },
        updateMeeting: async(_, { ID, name, desc, startTime, endTime, room, host, users }, { dataSources: ds }) => {
            const res = await ds.mrmsAPI('Meeting.Update', {info: {ID, name, desc, startTime, endTime, room, host, users}})
            return res
        },
    },

    Device: {
        owner: async({ owner: ID }, _, { dataSources: ds }) => {
            if (!ID) {
                return
            }
            const { Data: data = {}} = await ds.selfAPI('User.BatchFind', {'UUID': [ID]})
            return data[ID]
        },
    },
    Room: {
        devices: async({ devices: IDs }, _, { dataSources: ds }) => {
            if (!IDs) {
                return
            }
            const { data = {}} = await ds.mrmsAPI('Device.BatchFind', {'UUID': IDs})
            return IDs.map(ID => data[ID])
        },
    },
    Meeting: {
        room: async({ room: ID }, _, { dataSources: ds }) => {
            if (!ID) {
                return
            }
            const { data = {}} = await ds.mrmsAPI('Room.BatchFind', {'UUID': [ID]})
            return data[ID]
        },
        host: async({ host: ID }, _, { dataSources: ds }) => {
            if (!ID) {
                return
            }
            const { Data: data = {}} = await ds.selfAPI('User.BatchFind', {'UUID': [ID]})
            return data[ID]
        },
        users: async({ users: IDs }, _, { dataSources: ds }) => {
            if (!IDs) {
                return
            }
            const { Data: data = {}} = await ds.selfAPI('User.BatchFind', {'UUID': IDs})
            return IDs.map(ID => data[ID])
        },
    },
}

// TODO: use data-loader
