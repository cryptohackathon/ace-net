#!/usr/bin/env node
import { app } from "./app";
import * as WebSocket from 'ws';

const port = process.env.SERVER_PORT

export const wsServer = new WebSocket.Server({ noServer: true });
wsServer.on('connection', (socket: WebSocket) => {
    socket.on('message', (message: string) => console.log(message));
    socket.send(JSON.stringify({
        message: 'Hi there, I am a WebSocket server'
    }));
});

const server = app.listen(port, () =>
    // tslint:disable-next-line:no-console
    console.log(`server started listening at http://localhost:${ port }`)
);

server.on('upgrade', (request, socket, head) => {
    wsServer.handleUpgrade(request, socket, head, sock => {
        wsServer.emit('connection', sock, request);
    });
});
