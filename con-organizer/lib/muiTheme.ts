'use client';

import { createTheme, ThemeOptions } from '@mui/material';

const muiLightTheme: ThemeOptions = {
        palette: {
            primary: {
                light: '#fff',
                main: '#222',
                dark: '#000',
                contrastText: '#fff',
            },
            secondary: {
                light: '#fff',
                main: '#f44336',
                dark: '#000',
                contrastText: '#000',
            },
        },
    }
;

export const muiLight = createTheme(muiLightTheme);