const { gql } = require('apollo-server');

const typeDefs = gql`
scalar Timestamp

type Query {
    device(ID: ID!): Device
    devices(
        """
        Default = 5
        """
        size: Int,
        page: Int,
        word: String
    ): DeviceList

    room(ID: ID!): Room
    rooms(
        """
        Default = 5
        """
        size: Int,
        page: Int,
        word: String
    ): RoomList

    meeting(ID: ID!): Meeting
    meetings(
        """
        Default = 5
        """
        size: Int,
        page: Int,
        word: String
    ): MeetingList
}

type Mutation {
    createDevice(
        name: String!,
        type: String!,
	    owner: String
    ): Device
    deleteDevice(ID: String!): Boolean!
    updateDevice(
        ID: String!,
        name: String,
        type: String,
	    owner: String
    ): Device

    createRoom(
        name: String!,
        addr: String!,
	    devices: [String]
    ): Room
    deleteRoom(ID: String!): Boolean!
    updateRoom(
        ID: String!,
        name: String,
        addr: String,
	    devices: [String]
    ): Room

    createMeeting(
        name: String!,
        desc: String!,
        startTime: Timestamp!,
        endTime: Timestamp!,
        room: String!,
        host: String!,
	    users: [String]
    ): Meeting
    deleteMeeting(ID: String!): Boolean!
    updateMeeting(
        ID: String!,
        name: String,
        desc: String,
        startTime: Timestamp,
        endTime: Timestamp,
        room: String,
        host: String,
	    users: [String]
    ): Meeting
}

# ===

type Device {
    ID: ID!
    name: String!
    type: String!
	owner: User
}

type Room {
    ID: ID!
    name: String!
    addr: String!
	devices: [Device]
}

type Meeting {
    ID: ID!
    name: String!
    desc: String!
    startTime: Timestamp!
    endTime: Timestamp!

	room: Room
    host: User
    users: [User]
}

# ===

type DeviceList {
  total: Int!
  list: [Device]
}

type RoomList {
  total: Int!
  list: [Room]
}

type MeetingList {
  total: Int!
  list: [Meeting]
}

# ===

type User {
    # public
	UUID: String!
	nick: String!
	avatar: String!
	motto: String!
	homepage: String!

    # private
	name: String
	phone: String
	email: String
}
`;

module.exports = typeDefs;
