'use client';
import { ReactNode } from 'react';
import { ThemeProvider } from '@mui/material';
import { AuthProvider } from '@/components/AuthProvider';
import { muiDark } from '@/lib/muiTheme';

export default function Index({ children }: { children: ReactNode }) {
    return (
        <ThemeProvider theme={muiDark}>
            <AuthProvider>{children}</AuthProvider>
        </ThemeProvider>
    );
}
