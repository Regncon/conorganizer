import { AppBar, Paper, Tab, Tabs, Toolbar, Typography } from '@mui/material';
import RoomMap from './RoomMap';
const Rooms = async () => {
    const pool = 'Lørdag Morgen';
    return (
        <Paper
            sx={{
                width: 'calc(2901px + 2rem)',
                height: 'calc(2073px + 7rem)',
                position: 'absolute',
                left: '0',
                top: '60px',
                padding: '1rem',
                margin: '1rem',
            }}
        >
            {' '}
            <AppBar position="fixed" sx={{ paddingTop: '60px' }}>
                <Toolbar>
                    <Typography variant="h1">Romfordeling </Typography>
                    <Tabs value={1}  aria-label="basic tabs example">
                        <Tab label="Fredag Kveld" />
                        <Tab label="Lørdag Morgen" />
                        <Tab label="Lørdag Kveld" />
                        <Tab label="Søndag Morgen" />
                    </Tabs>
                </Toolbar>
            </AppBar>
            <Toolbar />
            <RoomMap pool="Lørdag Morgen" />
        </Paper>
    );
};

export default Rooms;
