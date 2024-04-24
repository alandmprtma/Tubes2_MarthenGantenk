import { Inter } from "next/font/google";
import "./globals.css";
import StarsCanvas from "@/components/StarBackground";

const inter = Inter({ subsets: ["latin"] });

export const metadata = {
  title: "Lemanspedia - WikiRace Program",
  description: "WikiRace Program with BFS and IDS Algorithms",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body className={inter.className}>
      <StarsCanvas/>
        {children}
        </body>
      
    </html>
  );
}
