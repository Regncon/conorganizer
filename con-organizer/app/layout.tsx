import { ReactNode } from 'react';
import { Inter } from 'next/font/google';
import Index from './Index';
import './globals.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
    title: 'Regncon program 2023',
    description: 'Regncon program og puljepåmelding 2023',
};
export default function RootLayout({ children }: { children: ReactNode }) {
    return (
        <html lang="en">
            <body className={[inter.className].join(' ').trim()}>
                {<Index>{children}</Index>}
            </body>
        </html>
    );
}
