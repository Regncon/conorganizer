import { Box, Container, CssBaseline, Paper, ThemeProvider } from '@mui/material';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v14-appRouter';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { muiDark } from '$lib/muiTheme';
import styles from './page.module.scss';
import './global.scss';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
	title: 'Regncon program 2024',
	description: 'Regncon program og puljepåmelding 2024',
};

export default function RootLayout({
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
							<Box component={Paper} className={styles['main-test']} elevation={1}>
								{children}
							</Box>
						</Container>
					</ThemeProvider>
				</AppRouterCacheProvider>
			</body>
		</html>
	);
}
