import { Box } from '@mui/material';
import MyParticipant from './UI/MyParticipantCard';
import { AssignParticipantByEmail } from './lib/actions/actions';
import MyParticipantsHeader from './MyParticipantsHeader';

const MyParticipants = async () => {
    const participants = await AssignParticipantByEmail();
    return (
        <Box sx={{ display: 'grid', height: 'var(--centering-height)', placeContent: 'center' }}>
            <Box>
                <MyParticipantsHeader participants={participants} />
                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                        gap: '2rem',
                    }}
                >
                    {participants?.map((participant) => (
                        <MyParticipant key={participant.id} participant={participant} />
                    ))}
                </Box>
            </Box>
        </Box>
    );
};

export default MyParticipants;
