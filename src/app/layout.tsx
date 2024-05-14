import { Box, Card, CardContent, Container, CssBaseline, ThemeProvider, Typography } from '@mui/material';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v14-appRouter';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { muiDark } from '$lib/muiTheme';
import styles from './page.module.scss';
import './global.scss';
import Link from 'next/link';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import BackButton from './BackButton';
import LogOutButton from './LogOutButton';

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
    const { auth } = await getAuthorizedAuth();

    return (
        <html lang="en">
            <body className={inter.className}>
                <AppRouterCacheProvider options={{ key: 'mui-theme' }}>
                    <ThemeProvider theme={muiDark}>
                        <CssBaseline enableColorScheme />
                        <Container component={'main'} maxWidth="xl">
                            <Box sx={{ display: 'flex', placeContent: 'space-between', paddingTop: '1rem' }}>
                                <BackButton />
                                <LogOutButton />
                            </Box>
                            <Box className={styles['main-test']}>
                                {auth?.currentUser?.uid ? null : (
                                    <Card sx={{ marginTop: '1rem' }}>
                                        <CardContent>
                                            <Typography variant="h1">
                                                For og lage arrangementer må du ha en bruker trykk på{' '}
                                                <Link href="/login">logginn</Link> Eller
                                                <Link href="/register"> registrer </Link>
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                )}

                                {children}
                            </Box>
                        </Container>
                    </ThemeProvider>
                </AppRouterCacheProvider>
            </body>
        </html>
    );
}
