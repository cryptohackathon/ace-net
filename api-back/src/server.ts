#!/usr/bin/env node
import { app } from "./app";

const port = process.env.SERVER_PORT


app.listen(port, () =>
    // tslint:disable-next-line:no-console
  console.log(`server started listening at http://localhost:${port}`)
);
