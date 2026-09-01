import type { Metadata } from "next";
import "./styles.css";
export const metadata: Metadata = { title: "GhanaCalendar — Working days, without guesswork", description: "Evidence-backed Ghana public holidays and deterministic working-day calculations.", metadataBase: new URL("https://calendar.digitalghana.dev") };
export default function Layout({children}:{children:React.ReactNode}){return <html lang="en"><body>{children}</body></html>}
