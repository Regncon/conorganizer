import { Box, Container, CssBaseline, Paper, ThemeProvider, Typography } from '@mui/material';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v14-appRouter';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { muiDark } from '$lib/muiTheme';
import styles from './page.module.scss';
import './global.scss';
import Link from 'next/link';
import { headers } from 'next/headers';
import { firebaseAuth } from '$lib/firebase/firebase';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';

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
							<Box className={styles['main-test']}>
								{auth?.currentUser?.uid ? null : (
									<Typography variant="h1">
										For og lage arrangementer må du ha en bruker trykk på{' '}
										<Link href="/login">logg inn</Link> Eller
										<Link href="/register"> registrer </Link>
									</Typography>
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
