'use client';

import { createTheme, ThemeOptions } from '@mui/material';

const muiLightTheme: ThemeOptions = {
    palette: {
        primary: {
            light: '#a1887f',
            main: '#3e2723',
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
    },
    components: {
        MuiButton: {
            styleOverrides: {
                root: ({ ownerState }) => ({
                    ...(ownerState.variant === 'contained' &&
                        ownerState.color === 'primary' && {
                           backgroundColor: 'primary.dark',
                            color: '#fff',
                        }),
                }),
            },
        },
    },
};
export const muiLight = createTheme(muiLightTheme);