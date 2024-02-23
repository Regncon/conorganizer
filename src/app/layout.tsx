import type { Metadata } from 'next';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v14-appRouter';
import { Inter } from 'next/font/google';
import { CssBaseline, ThemeProvider } from '@mui/material';
import { muiDark } from '@/lib/muiTheme';

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
						<CssBaseline />
						{children}
					</ThemeProvider>
				</AppRouterCacheProvider>
			</body>
		</html>
	);
}
