import { Box, Card, CardActionArea, Typography } from '@mui/material';
import ParticipantCard from './components/ParticipantCard';
import { Participant } from '$lib/types';
import AddIcon from '@mui/icons-material/Add';
import Link from 'next/link';

const participants = () => {
    const tmpParticipants: Participant[] = [];
    return (
        <Box>
            <Typography variant="h1">Participants</Typography>
            <Typography variant="h2">Under utvikilig. Leker med ekte data, Ikke trykk på ting. </Typography>
            <Box
                sx={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                    gap: '1rem',
                }}
            >
                <Link href="/admin/dashboard/participants/add" prefetch style={{ textDecoration: 'none' }}>
                    <Card
                        sx={{
                            minWidth: '306px',
                            height: '4rem',
                            width: '100%',
                            backgroundSize: 'cover',
                            borderRadius: '1.75rem',
                            opacity: 0.7,
                        }}
                    >
                        <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                            <AddIcon sx={{ fontSize: '4rem' }} />
                        </Box>
                    </Card>
                </Link>
                {tmpParticipants.map((participant) => (
                    <ParticipantCard key={participant.id} participant={participant} />
                ))}
            </Box>
        </Box>
    );
};
export default participants;
