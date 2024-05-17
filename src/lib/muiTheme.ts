'use client';
import { createTheme, ThemeOptions } from '@mui/material';
import { EB_Garamond, Raleway } from 'next/font/google';

const Garamond = EB_Garamond({
    weight: ['400', '700'],
    subsets: ['latin'],
    style: ['normal', 'italic'],
    display: 'swap',
});
const raleway = Raleway({
    weight: ['500', '700'],
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
            contrastText: '#FF7979',
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
            fontSize: '2.3rem',
        },
        h3: {
            fontWeight: '700',
            lineHeight: '1',
            fontSize: '1.8rem',
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
            fontFamily: raleway.style.fontFamily,
            fontSize: '1rem',
        },
    },
};

export const muiDark = createTheme(muiDarkTheme);
