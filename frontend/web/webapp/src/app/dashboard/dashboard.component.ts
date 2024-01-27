import { Component, inject, OnInit } from '@angular/core';
import { Breakpoints, BreakpointObserver } from '@angular/cdk/layout';
import { map } from 'rxjs/operators';
import { HttpClient } from "@angular/common/http";
import { MessageData, SseService } from '../sse-service';
import { MatSnackBar } from '@angular/material/snack-bar';


interface Log {
  message: string,
  level: string,
  date: string
}

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit {
  private breakpointObserver = inject(BreakpointObserver);

  dataSource = [];
  displayedColumns: string[] = ['message', 'level', 'date'];


  /** Based on the screen size, switch from standard to one column per row */
  cards = this.breakpointObserver.observe(Breakpoints.Handset).pipe(
    map(({ matches }) => {
      if (matches) {
        return [
          { title: 'Card 1', cols: 1, rows: 4 },
          // { title: 'Card 2', cols: 1, rows: 1 },
          // { title: 'Card 3', cols: 1, rows: 1 },
          // { title: 'Card 4', cols: 1, rows: 1 }
        ];
      }

      return [
        { title: 'Card 1', cols: 2, rows: 2 },
        // { title: 'Card 2', cols: 2, rows: 1 },
        // { title: 'Card 3', cols: 1, rows: 2 },
        // { title: 'Card 4', cols: 1, rows: 1 }
      ];
    })
  );

  constructor(private http: HttpClient, private sseService: SseService, private snackBar: MatSnackBar) {
  }

  ngOnInit() {
    this.http.get<any>("/logs").subscribe((log) => {
      this.dataSource = log.message
    });

    // this.messageService.add({ severity: 'succes', summary: 'Warn', detail: 'Message Content', life: 20000 });
    console.log("Created event source")
    this.sseService.createEventSource().subscribe(
      (e: MessageData) => {
        // this.messageService.add({ severity: 'warn', summary: 'Warn', detail: 'Message Content', life: 20000 });
        console.log('Message received: ' + e.message);
        this.snackBar.open("Detected a new error: " + e.message, "OK", {
          duration: 5000
        });
      }
    );
  }
}
