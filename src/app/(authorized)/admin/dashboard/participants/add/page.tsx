import {
    EventTicket,
    GetTicketsFromCheckIn,
} from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import TicketList from './TicketList';

const AddParticipant = async () => {
    const tickets: EventTicket[] | undefined = await GetTicketsFromCheckIn();
    return <TicketList tickets={tickets || []} />;
};

export default AddParticipant;
