// import { Container } from "typescript-ioc";

// const iocContainer = Container

// export { iocContainer };


// src/ioc.ts
import { IocContainer, IocContainerFactory } from "@tsoa/runtime";
import { Container } from "typescript-ioc";
// import { Container } from "di-package";

// Assign a container to `iocContainer`.
// const iocContainer = new Container();

// Or assign a function with to `iocContainer`.
const iocContainer: IocContainerFactory = (
  request: Request
) => {
  const container = Container;
//   container.bind(request);
  return container;
};

// export according to convention
export { iocContainer };
