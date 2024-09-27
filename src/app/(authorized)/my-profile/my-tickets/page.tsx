import { redirect } from 'next/navigation';
import MyParticipants from './components/MyParticipants';
import { AssignParticipantByEmail } from './components/lib/actions/actions';

const MyTickets = async () => {
    const participants = await AssignParticipantByEmail();

    if (participants?.length === 0 || participants === undefined) {
        redirect('/my-profile/my-tickets/not-found');
    }

    return <MyParticipants />;
};

export default MyTickets;
