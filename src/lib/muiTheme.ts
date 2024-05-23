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
            main: '#ff7c7c',
            contrastText: '#000000',
        },
        secondary: {
            main: '#ffd082',
            contrastText: '#000000',
        },
    background: {
        default: '#21243d',
        paper: '#31354e',
      },
      divider: '#88E1F2',
    },
    typography: {

        h1: {
            fontWeight: '700',
            fontFamily: Garamond.style.fontFamily,
            fontSize: '36px',
            marginBlock: '1rem',
        },
        h2: {
            fontWeight: '700',
            fontFamily: Garamond.style.fontFamily,
            fontSize: '24px',
            marginBlock: '1rem',
        },
        h3: {
            fontWeight: '700',
            lineHeight: '1',
            fontSize: '22px',
            // textShadow: "0 0 1em black",
            fontFamily: Garamond.style.fontFamily,
            marginBlock: '1rem',
        },
        h4: {
            color: '#ddd',
            fontWeight: '400',
            fontStyle: 'italic',
            fontSize: '20px',
            // textShadow: "0 0 .7em black",
            fontFamily: Garamond.style.fontFamily,
            marginBlock: '1rem',
        },
        h5: {
            fontWeight: 'bold',
            fontSize:"16px" 
        },
        h6: {
            fontWeight: 'bold',
            fontSize:"14px"
        },
        body1: {
            fontFamily: raleway.style.fontFamily,
            fontSize: '16px',
        },
    },
};

export const muiDark = createTheme(muiDarkTheme);
