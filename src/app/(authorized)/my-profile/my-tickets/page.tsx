import { Box, Button } from '@mui/material';
import ConfirmOrBuy from './MyTicketsComponents/ConfirmOrBuy';
import TicketNotFound from './MyTicketsComponents/TicketNotFound';
import Tickets from './MyTicketsComponents/Tickets';

const MyTickets = async () => {
    return (
        <Box sx={{ display: 'flex', paddingLeft: '2rem', gap: '1rem' }}>
            <ConfirmOrBuy />
            <TicketNotFound />
            <Tickets />

            <Button variant="contained" color="primary" href="/my-profile">
                Go back to my profile
            </Button>
        </Box>
    );
};

export default MyTickets;
