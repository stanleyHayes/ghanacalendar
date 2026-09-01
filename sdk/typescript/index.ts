export type GhanaHoliday={id:string;name:string;date:string|null;observedDate?:string;kind:string;status:"confirmed"|"pending";sourceId:string};
export type WorkingDayResult={date:string;timezone:"Africa/Accra";workingDay:boolean;holiday:GhanaHoliday|null};

export class GhanaCalendarClient{
  constructor(private readonly baseUrl="https://api-calendar.digitalghana.dev"){}
  private async get<T>(path:string):Promise<T>{const response=await fetch(`${this.baseUrl}${path}`);if(!response.ok)throw new Error(`GhanaCalendar request failed (${response.status})`);return response.json() as Promise<T>}
  holidays(year:number,includePending=false){return this.get<{version:string;timezone:string;holidays:GhanaHoliday[]}>(`/v1/holidays?year=${year}&includePending=${includePending}`)}
  isWorkingDay(date:string){return this.get<WorkingDayResult>(`/v1/working-day?date=${encodeURIComponent(date)}`)}
  nextWorkingDay(date:string){return this.get<{date:string;result:string;timezone:string}>(`/v1/next-working-day?date=${encodeURIComponent(date)}`)}
  previousWorkingDay(date:string){return this.get<{date:string;result:string;timezone:string}>(`/v1/previous-working-day?date=${encodeURIComponent(date)}`)}
  addWorkingDays(date:string,count:number){return this.get<{date:string;count:number;result:string;timezone:string}>(`/v1/add-working-days?date=${encodeURIComponent(date)}&count=${count}`)}
}
