import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { catchError, filter, map, shareReplay, take, tap } from 'rxjs/operators';
import { WsDataService } from '../services/ws-data.service';
import { ChartDataSets, ChartOptions, ChartType } from 'chart.js';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';


interface ApiDefaultResponse {
  /**
   * Optional details for error responses. The type depend on status.
   */
  errorDetails?: any;
  /**
   * Simple message to explain client developers the reason for error.
   */
  errorMessage?: string;
  /**
   * Response status. OK for successful reponses.
   */
  status: string;
}

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  constructor(
    private wsDataService: WsDataService,
    private httpClient: HttpClient
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

  data$!: Observable<any>
  slotLabels$!: Observable<any>

  ngOnInit(): void {
    this.wsDataService.connect()
    this.data$ = this.wsDataService.socket$.pipe(
      filter(x => x.message == 'POOLS_SUMMARY'),
      map(x => x.data),
      shareReplay(1)
    )
    this.slotLabels$ = this.data$.pipe(
      filter(data => data && data.length > 0),
      map(data => data[0].slotLabels)
    )
    // console.log(this.wsDataService.messages$)
  }

  histogramValue(histogram: number[], i: number) {
    return histogram && histogram.length > i ? histogram[i] : ""
  }

  chartData(histogram: number[]) {
    let data = [
      { data: histogram, label: 'Histogram' }
    ] as ChartDataSets[]
    // console.log(data)
    return data

  }

  barChartOptions: ChartOptions = {
    responsive: false,
    animation: {
      duration: 0
    }
  };
  barChartType: ChartType = 'bar';
  barChartLegend = false;
  barChartPlugins = [];

  statusColor(status: string) {
    switch (status) {
      case 'REGISTRATION':
        return 'gray'
      case 'ENCRYPTION':
        return 'orange'
      case 'FINALIZED':
        return 'lightgreen'
      case 'CALCULATED':
        return 'green'
      case 'EXPIRED':
        return 'red'
      default:
        return 'black'
    }
  }

  async reset() {
    console.log("RESET")
    await this.httpClient.post<ApiDefaultResponse>(`${ environment.basePath }/ace/reset`,
      null,
      {
        // params: queryParameters,
        // withCredentials: this.configuration.withCredentials,
        // headers: headers,
        // observe: observe,
        // reportProgress: reportProgress
      }
    ).pipe(take(1)).toPromise()
  }
}
