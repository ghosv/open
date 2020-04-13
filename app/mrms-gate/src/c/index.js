"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var request = require("request-promise-native");
var WebSocket = require("ws");
var url = require("url");
var defaultLocal = "http://localhost:8080/client";
var defaultLive = "https://api.micro.mu/client";
var Stream = /** @class */ (function () {
    function Stream(conn, service, endpoint) {
        this.conn = conn;
        this.service = service;
        this.endpoint = endpoint;
    }
    Stream.prototype.send = function (msg) {
        var _this = this;
        return new Promise(function (resolve, reject) { return __awaiter(_this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                this.conn.send(marshalRequest(this.service, this.endpoint, msg));
                return [2 /*return*/];
            });
        }); });
    };
    // this probably should use observables or something more modern
    Stream.prototype.recv = function (cb) {
        this.conn.on("message", function (m) {
            cb(unmarshalResponse(m));
        });
    };
    return Stream;
}());
exports.Stream = Stream;
var Client = /** @class */ (function () {
    function Client(options) {
        this.options = {
            address: defaultLive
        };
        if (options) {
            this.options = options;
        }
        if (options && options.local) {
            this.options.address = defaultLocal;
        }
    }
    // Call enables you to access any endpoint of any service on Micro
    Client.prototype.call = function (service, endpoint, req) {
        var _this = this;
        return new Promise(function (resolve, reject) { return __awaiter(_this, void 0, void 0, function () {
            var serviceReq, options, result, e_1;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        _a.trys.push([0, 2, , 3]);
                        // example curl: curl -XPOST -d '{"service": "go.micro.srv.greeter", "endpoint": "Say.Hello"}'
                        //  -H 'Content-Type: application/json' http://localhost:8080/client {"body":"eyJtc2ciOiJIZWxsbyAifQ=="}
                        if (!req) {
                            req = {};
                        }
                        serviceReq = {
                            service: service,
                            endpoint: endpoint,
                            // !
                            request: req
                        };
                        options = {
                            method: "POST",
                            json: true,
                            headers: {
                                micro_token: this.options.token
                            },
                            body: serviceReq
                        };
                        options.uri = this.options.address;
                        return [4 /*yield*/, request.post(this.options.address, options)];
                    case 1:
                        result = _a.sent();
                        // !
                        resolve(result);
                        return [3 /*break*/, 3];
                    case 2:
                        e_1 = _a.sent();
                        reject(e_1);
                        return [3 /*break*/, 3];
                    case 3: return [2 /*return*/];
                }
            });
        }); });
    };
    Client.prototype.stream = function (service, endpoint, msg) {
        var _this = this;
        return new Promise(function (resolve, reject) {
            try {
                var uri = url.parse(_this.options.address);
                // TODO: make optional
                uri.path = "/client/stream";
                uri.pathname = "/client/stream";
                uri.protocol = uri.protocol.replace("http", "ws");
                var conn_1 = new WebSocket(url.format(uri), {
                    perMessageDeflate: false
                });
                var data_1 = marshalRequest(service, endpoint, msg);
                conn_1.on("open", function open() {
                    conn_1.send(data_1);
                    var stream = new Stream(conn_1, service, endpoint);
                    resolve(stream);
                });
            }
            catch (e) {
                reject(e);
            }
        });
    };
    return Client;
}());
exports.Client = Client;
function marshalRequest(service, endpoint, v) {
    var jsonBody = JSON.stringify(v);
    return JSON.stringify({
        service: service,
        endpoint: endpoint,
        body: Buffer.from(jsonBody).toString("base64")
    });
}
function unmarshalResponse(msg) {
    var rsp = JSON.parse(msg);
    return Buffer.from(rsp.body, "base64").toString();
}
