import {
    EventTicket,
    GetTicketsFromCheckIn,
} from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import { Typography } from '@mui/material';
import TicketList from './components/TicketList';

const AddParticipant = async () => {
    const tickets: EventTicket[] | undefined = await GetTicketsFromCheckIn();

    return (
        <>
            <Typography variant="h1">Ligg til deltager</Typography>
            <TicketList tickets={tickets || []} />
        </>
    );
};

export default AddParticipant;
