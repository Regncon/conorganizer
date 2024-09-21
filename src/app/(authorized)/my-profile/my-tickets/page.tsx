import { Box } from '@mui/material';
import TicketNotFound from './components/TicketNotFound';
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

    return (
        <Box sx={{ display: 'flex', paddingLeft: '2rem', gap: '1rem' }}>
            {tickets?.length === 0 ?
                <TicketNotFound />
                : <Tickets tickets={tickets as EventTicket[]} />}
        </Box>
    );
};

export default MyTickets;
