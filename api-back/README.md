# ACE api-back

## Introduction

This is the backend for ACE* project. It implements the pool management system in four phases:
- client registration to pool,
- clients key sharing,
- submission of encrypted data and keys to apply an inner product to the data,
- inner product calculation.

## Getting started 

```
npm install
npm run routes
npm run dev
```

The project is based on [TSOA](https://tsoa-community.github.io/docs/). 


## API documentation

- [static Swagger API documentation](https://cryptohackathon.github.io/ace-net/),
- check `http://localhost:9500/api-doc/` for testing.


## Websocket demo

Use Chrome. Google to Chrome Web Store. Install *Simple Websocket Client* extension.
While server is running (`npm run dev`) do the following.
Connect with *Simple Websocket Client* to ws://localhost:9500  (open connection, see https://medium.com/factory-mind/websocket-node-js-express-step-by-step-using-typescript-725114ad5fe4).
Go to folder `../test-clients`. Run script `./run-clients` which runs 100 clients. 
There should be many messages describing communication with between clients and server.

