'use client';

import { createTheme, ThemeOptions } from '@mui/material';
import { EB_Garamond } from 'next/font/google';

const Garamond = EB_Garamond({
    weight: ['400', '700'],
    subsets: ['latin'],
    style: ['normal', 'italic'],
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
            fontSize: '2rem',
        },
        h3: {
            fontWeight: '700',
            lineHeight: '1',
            fontSize: '2.5rem',
            fontFamily: Garamond.style.fontFamily,
        },
        h4: {
            color: 'lightgrey',
            fontWeight: '400',
            fontStyle: 'italic',
            fontSize: '1.5rem',
            fontFamily: Garamond.style.fontFamily,
        },
    },
};

export const muiDark = createTheme(muiDarkTheme);
