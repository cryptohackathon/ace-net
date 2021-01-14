import { Injectable } from '@angular/core';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { environment } from '../../environments/environment';
import { catchError, tap, switchAll } from 'rxjs/operators';
import { EMPTY, Subject } from 'rxjs';
export const WS_ENDPOINT = environment.wsEndpoint;

// export class Message {
//   constructor(
//     public sender: string,
//     public content: string,
//     public isBroadcast = false,
//   ) { }
// }

interface WSMessage {
  message: string;
}

@Injectable({
  providedIn: 'root'
})
export class WsDataService {
  public socket$!: WebSocketSubject<any>;
  private messagesSubject$ = new Subject<any>();
  public messages$ = this.messagesSubject$.pipe(
    switchAll(),
    catchError(e => { throw e })
  );

  public connect(): void {
    console.log("TRYING TO CONNECT")
    if (!this.socket$ || this.socket$.closed) {
      console.log("INITIALIZING WS")
      this.socket$ = this.getNewWebSocket() as WebSocketSubject<any>;
      const messages = this.socket$.pipe(
        tap({
          error: error => console.log(error),
        }), catchError(_ => EMPTY));
      this.messagesSubject$.next(messages);
    }
  }

  private getNewWebSocket() {
    return webSocket(WS_ENDPOINT);
  }
  sendMessage(msg: any) {
    this.socket$.next(msg);
  }
  close() {
    this.socket$.complete();
  }
}
