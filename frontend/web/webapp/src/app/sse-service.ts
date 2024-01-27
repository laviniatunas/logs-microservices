import { Observable } from "rxjs-observable";

export interface MessageData {
    message: string;
}

export class SseService {

    constructor() { }

    createEventSource(): Observable<MessageData> {
        const eventSource = new EventSource("http://localhost:8000/stream");

        return new Observable(observer => {
            eventSource.onmessage = event => {
                const messageData: MessageData = JSON.parse(event.data);
                observer.next(messageData);
            };
        });
    }
}
