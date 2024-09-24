import { Box, Typography } from '@mui/material';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Participant } from '$lib/types';
import MyParticipant from './UI/MyParticipantCard';

type Props = { participants: Participant[] | undefined };

const Tickets = async ({ participants }: Props) => {
    const { user } = await getAuthorizedAuth();
    const verifiedEmail = user?.emailVerified ?? false;
    const verifiedCheckIn = true;

    if (verifiedEmail && verifiedCheckIn) {
        return (
            <Box sx={{ display: 'grid', height: 'var(--centering-height)', placeContent: 'center' }}>
                <Box>
                    <Typography>En smart hjelpetekst skrevet av en som ikke er meg eller dyslektiker</Typography>
                    <Typography variant="h1">Mine billetter</Typography>
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
    }
    return null;
};

export default Tickets;