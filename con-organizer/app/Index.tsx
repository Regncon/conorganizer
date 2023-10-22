'use client';
import { ReactNode } from 'react';
import { CssBaseline, ThemeProvider } from '@mui/material';
import AppHeader from '@/components/UI/AppHeader';
import { AuthProvider } from '@/components/AuthProvider/AuthProvider';
import { muiDark } from '@/lib/muiTheme';

export default function Index({ children }: { children: ReactNode }) {
    return (
        <ThemeProvider theme={muiDark}>
            <CssBaseline enableColorScheme />
            <AppHeader />
            <AuthProvider>{children}</AuthProvider>
        </ThemeProvider>
    );
}
