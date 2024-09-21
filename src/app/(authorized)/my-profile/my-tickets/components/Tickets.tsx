import { Box, Paper, Typography } from '@mui/material';
import Ticket from './UI/Ticket';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { EventTicket } from './lib/actions/actions';

type Props = { tickets: EventTicket[] | undefined };

const Tickets = async ({ tickets }: Props) => {
    const { user } = await getAuthorizedAuth();
    const verifiedEmail = user?.emailVerified ?? false;
    const verifiedCheckIn = true;
    const testData: EventTicket[] = [
        {
            id: 1,
            category: 'Test',
            category_id: 1,
            crm: {
                first_name: 'Test',
                last_name: 'Testesen',
                id: 1,
                email: 'test@test.com',
                born: '1990-01-01',
            },
            order_id: 1,
        },
        {
            id: 2,
            category: 'Test',
            category_id: 1,
            crm: {
                first_name: 'Test',
                last_name: 'Testesen',
                id: 1,
                email: 'test@test.com',
                born: '1990-01-01',
            },
            order_id: 1,
        },
    ];
    if (verifiedEmail && verifiedCheckIn) {
        return (
            <Box sx={{ display: 'grid', height: 'var(--centering-height)', placeContent: 'center' }}>
                <Box>
                    <Typography>En smart hjelpetekst skrevet av en som ikke er meg eller dyslektiker</Typography>
                    <Typography variant="h1">My Tickets</Typography>
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                            gap: '2rem',
                        }}
                    >
                        {testData?.map((ticket) => <Ticket key={ticket.id} ticket={ticket} />)}
                    </Box>
                </Box>
            </Box>
        );
    }
    return null;
};

export default Tickets;
