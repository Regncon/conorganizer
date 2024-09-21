import Tickets from './components/Tickets';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { redirect } from 'next/navigation';
import { GetTicketsByEmail } from './components/lib/actions/actions';

const MyTickets = async () => {
    const { user } = await getAuthorizedAuth();
    if (user?.emailVerified === false) {
        redirect('/my-profile/my-tickets/confirm');
    }
    const tickets = await GetTicketsByEmail(user?.email);

    if (tickets?.length === 0 || tickets === undefined) {
        redirect('/my-profile/my-tickets/not-found');
    }

    return <Tickets tickets={tickets} />;
};

export default MyTickets;
