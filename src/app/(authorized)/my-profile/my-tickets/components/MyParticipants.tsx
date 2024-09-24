import { Box, Typography } from '@mui/material';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Participant } from '$lib/types';
import MyParticipant from './UI/MyParticipantCard';
import { AssignParticipantByEmail } from './lib/actions/actions';
import MyParticipantsHeader from './MyParticipantsHeader';

type Props = { participants: Participant[] | undefined };

const MyParticipants = async ({ participants }: Props) => {
    const { user } = await getAuthorizedAuth();
    const verifiedEmail = user?.emailVerified ?? false;
    const verifiedCheckIn = true;

    if (verifiedEmail && verifiedCheckIn) {
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
    }
    return null;
};

export default MyParticipants;
