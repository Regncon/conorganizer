import { Box } from '@mui/material';
import { Typography } from '@/lib/mui';

const AppHeader = () => {
    return (
        <Box sx={{ p: '2em', margin: '0 auto', maxWidth: '1080px' }}>
            <header className="AppHeader">
                <img src="/image/regnconlogony.png" alt="Regncondragen for 2023" className="regnconLogo" />
                <Typography variant="h1" color="white">
                    Regncon XXXI
                </Typography>
            </header>
        </Box>
    );
};

export default AppHeader;
