import { createTheme, ThemeProvider } from '@mui/material/styles';

interface ThemeProps {
    children: React.ReactNode;
}

export const Theme: React.FC<ThemeProps> = ({ children }) => {
    const tema = createTheme({
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
    });

    return <ThemeProvider theme={tema}>{children}</ThemeProvider>;
};
