import { Box, Typography } from '@mui/material';
import { getAllParticipants } from '$app/(public)/components/lib/serverAction';
import ParticipantsList from './components/ParticipantsList';

const participants = async () => {
    const tmpParticipants = await getAllParticipants();
    return (
        <Box>
            <Typography variant="h1">Participants</Typography>
            <Typography variant="h2">Under utvikilig. Leker med ekte data, Ikke trykk på ting. </Typography>
            <ParticipantsList participants={tmpParticipants} />
        </Box>
    );
};
export default participants;
