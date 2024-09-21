import { Box } from '@mui/material';
import TicketNotFound from './not-found/TicketNotFound';
import Tickets from './components/Tickets';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { redirect } from 'next/navigation';
import { EventTicket, GetTicketsByEmail } from './components/lib/actions/actions';

const MyTickets = async () => {
    const { user } = await getAuthorizedAuth();
    if (user?.emailVerified === false) {
        redirect('/my-profile/my-tickets/confirm');
    }
    const tickets = await GetTicketsByEmail(user?.email);

    if (tickets?.length === 0 || tickets === undefined) {
        redirect('/my-profile/my-tickets/not-found');
    }

    return (
        <Box sx={{ display: 'flex', paddingLeft: '2rem', gap: '1rem' }}>
            <Tickets tickets={tickets} />
        </Box>
    );
};

export default MyTickets;
