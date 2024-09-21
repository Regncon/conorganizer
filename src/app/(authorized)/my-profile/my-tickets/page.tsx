import { Box, Button } from '@mui/material';
import TicketNotFound from './components/TicketNotFound';
import Tickets from './components/Tickets';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { use } from 'react';
import { redirect } from 'next/navigation';

const MyTickets = async () => {
    const { user } = await getAuthorizedAuth();
    if (user?.emailVerified === false) {
        redirect('/my-profile/my-tickets/confirm');
    }
    return (
        <Box sx={{ display: 'flex', paddingLeft: '2rem', gap: '1rem' }}>
            <TicketNotFound />
            <Tickets />

            <Button variant="contained" color="primary" href="/my-profile">
                Go back to my profile
            </Button>
        </Box>
    );
};

export default MyTickets;
