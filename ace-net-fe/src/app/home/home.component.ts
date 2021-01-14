import { Component, OnInit } from '@angular/core';
import { catchError, map, tap } from 'rxjs/operators';
import { WsDataService } from '../services/ws-data.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  constructor(
    private wsDataService: WsDataService
  ) {
  }

  liveData$ = this.wsDataService.messages$.pipe(
    map(rows => (rows as any).data),
    catchError(error => { throw error }),
    tap({
      error: error => console.log('[Live component] Error:', error),
      complete: () => console.log('[Live component] Connection Closed')
    }
    )
  );

  data$ = this.wsDataService.socket$
  ngOnInit(): void {
    this.wsDataService.connect()
    this.data$ = this.wsDataService.socket$
    // console.log(this.wsDataService.messages$)
  }

}
