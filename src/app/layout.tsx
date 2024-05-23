import { Box, Card, CardContent, Container, CssBaseline, ThemeProvider, Typography } from '@mui/material';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v14-appRouter';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { muiDark } from '$lib/muiTheme';
import styles from './page.module.scss';
import './global.scss';
import Link from 'next/link';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import BackButton from './(authorized)/BackButton';
import LogOutButton from './(authorized)/LogOutButton';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
    title: 'Regncon program 2024',
    description: 'Regncon program og puljepåmelding 2024',
};

export default async function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
            <body className={inter.className}>
                <AppRouterCacheProvider options={{ key: 'mui-theme' }}>
                    <ThemeProvider theme={muiDark}>
                        <CssBaseline enableColorScheme />
                        <Container component={'main'} maxWidth="xl">
                            <Box className={styles['main-test']}>{children}</Box>
                        </Container>
                    </ThemeProvider>
                </AppRouterCacheProvider>
            </body>
        </html>
    );
}
