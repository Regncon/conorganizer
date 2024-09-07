import { Paper, Typography } from '@mui/material';
import RoomMap from './RoomMap';
const Rooms = async () => {
    const pool = 'Lørdag Morgen';
    return (
        <Paper sx={{ width: 'calc(3868px + 2rem)', height: 'calc(2764px + 7rem)', padding: '1rem', margin: '1rem' }}>
            <Typography variant="h1">Romfordeling: {pool} </Typography>
            <RoomMap pool="Lørdag Morgen" />
        </Paper>
    );
};

export default Rooms;
