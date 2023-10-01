import { Box } from '@mui/material';
import { Typography } from '@/lib/mui';

const AppHeader = () => {
    return (
        <Box sx={{ p: '5em 0 4em 1em', margin: '0 auto', maxWidth: '900px' }}>
            <header className="AppHeader">
                <img
                    src="/image/regnconlogony.png"
                    alt="Regncondragen for 2023"
                    className="regnconLogo"
                    onClick={() => (window.location.href = `/`)}
                />
                <div>
                    <Typography variant="h1" color="white">
                        Regncon XXXI
                    </Typography>
                    <Typography variant="h4">Program</Typography>
                </div>
            </header>
        </Box>
    );
};

export default AppHeader;
