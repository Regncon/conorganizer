'use client';
import { ReactNode } from 'react';
import { ThemeProvider } from '@mui/material';
import AppHeader from '@/components/AppHeader';
import { AuthProvider } from '@/components/AuthProvider';
import { muiDark } from '@/lib/muiTheme';

export default function Index({ children }: { children: ReactNode }) {
    return (
        <ThemeProvider theme={muiDark}>
            <AppHeader />
            <AuthProvider>{children}</AuthProvider>
        </ThemeProvider>
    );
}
