import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { redirect } from 'next/navigation';
import MyParticipants from './components/MyParticipants';
import { AssignParticipantByEmail } from './components/lib/actions/actions';

const MyTickets = async () => {
    const { user } = await getAuthorizedAuth();
    if (user?.emailVerified === false) {
        redirect('/my-profile/my-tickets/confirm');
    }
    const participants = await AssignParticipantByEmail();

    if (participants?.length === 0 || participants === undefined) {
        redirect('/my-profile/my-tickets/not-found');
    }

    return <MyParticipants participants={participants} />;
};

export default MyTickets;
