import { Singleton, Factory } from "typescript-ioc";

@Singleton
@Factory(() => new EnergyEngine())
export class EnergyEngine {
    public async test() {
        console.log("CALLING")
        return {
            message: "Hello World 2"
        }
    }
}