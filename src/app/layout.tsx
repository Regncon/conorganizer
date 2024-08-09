import { Box, Container, CssBaseline, ThemeProvider, useMediaQuery, useTheme } from '@mui/material';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v14-appRouter';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { muiDark } from '$lib/muiTheme';
import styles from './page.module.scss';
import './global.scss';

// <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png"></link>
// <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
// <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
// <link rel="manifest" href="/site.webmanifest">
// <link rel="mask-icon" href="/safari-pinned-tab.svg" color="#88E1F2">
// <meta name="msapplication-TileColor" content="#ff7c7c">
// <meta name="theme-color" content="#000000">

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
                        <Container maxWidth="xl" disableGutters component={'main'}>
                            <Box className={styles['main-test']}>{children}</Box>
                        </Container>
                    </ThemeProvider>
                </AppRouterCacheProvider>
            </body>
        </html>
    );
}
