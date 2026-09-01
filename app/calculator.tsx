"use client";
import { useMemo, useState } from "react";

type Holiday={name:string;date:string|null;observedDate?:string;status:string};
function iso(value:Date){return value.toISOString().slice(0,10)}
export function Calculator({holidays}:{holidays:Holiday[]}){
  const [date,setDate]=useState("2026-12-25"); const [count,setCount]=useState(2);
  const holidayDates=useMemo(()=>new Set(holidays.flatMap(h=>[h.date,h.observedDate].filter(Boolean) as string[])),[holidays]);
  const result=useMemo(()=>{const cursor=new Date(`${date}T12:00:00Z`);let remaining=Math.abs(count),step=count<0?-1:1;while(remaining>0){cursor.setUTCDate(cursor.getUTCDate()+step);const day=cursor.getUTCDay(),value=iso(cursor);if(day!==0&&day!==6&&!holidayDates.has(value))remaining--}return iso(cursor)},[date,count,holidayDates]);
  const selected=holidays.find(h=>h.date===date||h.observedDate===date);const day=new Date(`${date}T12:00:00Z`).getUTCDay();const working=day!==0&&day!==6&&!holidayDates.has(date);
  return <div className="calculator"><div><label htmlFor="date">Start date</label><input id="date" type="date" min="2024-01-01" max="2026-12-31" value={date} onChange={e=>setDate(e.target.value)}/></div><div><label htmlFor="count">Working days to add</label><input id="count" type="number" min="-60" max="60" value={count} onChange={e=>setCount(Number(e.target.value))}/></div><div className="answer"><span>{working?"Working day":selected?.name??"Weekend"}</span><strong>{result}</strong><small>Result · Africa/Accra</small></div></div>
}
