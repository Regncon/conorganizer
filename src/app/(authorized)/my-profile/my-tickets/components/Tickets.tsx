import { Box, Paper, Typography } from '@mui/material';
import Ticket from './UI/Ticket';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { EventTicket } from './lib/actions/actions';

type Props = { tickets: EventTicket[] };

const Tickets = async ({ tickets }: Props) => {
    const { user } = await getAuthorizedAuth();
    const verifiedEmail = user?.emailVerified ?? false;
    const verifiedCheckIn = true;

    if (verifiedEmail && verifiedCheckIn) {
        return (
            <Box sx={{ display: 'grid', height: 'var(--centering-height)', placeContent: 'center' }}>
                <Typography>En smart hjelpetekst skrevet av en som ikke er meg eller dyskelktiker</Typography>
                <Typography variant="h1">My Tickets</Typography>
                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                        gap: '1rem',
                    }}
                >
                    {tickets.map((ticket) => (
                        <Ticket key={ticket.id} ticket={ticket} />
                    ))}
                </Box>
            </Box>
        );
    }
    return null;
};

export default Tickets;
