'use client';

import { createTheme, ThemeOptions } from '@mui/material';
import { EB_Garamond, Inter } from 'next/font/google';

const Garamond = EB_Garamond({
	weight: ['400', '700'],
	subsets: ['latin'],
	style: ['normal', 'italic'],
	display: 'swap',
});
const inter = Inter({
	weight: ['400', '700'],
	subsets: ['latin'],
	style: ['normal'],
	display: 'swap',
});

const muiDarkTheme: ThemeOptions = {
	palette: {
		mode: 'dark',
		primary: {
			light: '#e0cfc9',
			main: '#a1887f',
			dark: '#000',
			contrastText: '#fff',
		},
		secondary: {
			light: '#ffd54f',
			main: '#ff8f00',
			dark: '#000',
			contrastText: '#000',
		},
	},
	typography: {
		h6: {
			fontWeight: 'bold',
		},
		h1: {
			fontWeight: '700',
			fontFamily: Garamond.style.fontFamily,
			fontSize: '2.7rem',
		},
		h2: {
			fontWeight: '700',
			fontFamily: Garamond.style.fontFamily,
			color: 'grey',
			fontStyle: 'italic',
			fontSize: '2.5em',
			textShadow: '0 0 1em black',
		},
		h3: {
			fontWeight: '700',
			lineHeight: '1',
			fontSize: '2.2rem',
			// textShadow: "0 0 1em black",
			fontFamily: Garamond.style.fontFamily,
		},
		h4: {
			color: '#ddd',
			fontWeight: '400',
			fontStyle: 'italic',
			fontSize: '1.3rem',
			// textShadow: "0 0 .7em black",
			fontFamily: Garamond.style.fontFamily,
		},
		body1: {
			fontFamily: inter.style.fontFamily,
			fontSize: '1rem',
		},
	},
};

export const muiDark = createTheme(muiDarkTheme);
