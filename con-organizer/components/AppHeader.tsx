import * as React from 'react';
import { Box } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import useMediaQuery from '@mui/material/useMediaQuery';
import { Typography } from '@/lib/mui';

const AppHeader = () => {
    const theme = useTheme();
    const isSmallScreen = useMediaQuery(theme.breakpoints.down('sm'));
    console.log(isSmallScreen);

    return (
        <Box
            sx={
                isSmallScreen
                    ? { p: '1em 0 1em 1em', margin: '0 auto', maxWidth: '900px' }
                    : { p: '5em 0 4em 1em', margin: '0 auto', maxWidth: '600px' }
            }
        >
            <header className="AppHeader">
                <img
                    src="/image/regnconlogony.png"
                    alt="Regncondragen for 2023"
                    className="regnconLogo"
                    onClick={() => (window.location.href = `/`)}
                />
                <div>
                    <Typography variant={isSmallScreen ? 'h5' : 'h1'} color="white">
                        Regncon XXXI
                    </Typography>
                    <Typography variant="h4">Program</Typography>
                </div>
            </header>
        </Box>
    );
};

export default AppHeader;
