import { Participant } from '$lib/types';
import { User } from 'firebase/auth';
import { EventTicket } from './actions';

export const NewTickets = async (tickets: EventTicket[], participants: Participant[], user: User) => {
    console.log('Assigning new tickets to participants');
    // console.log('tickes', tickets, 'participants', participants, 'user', user);

    const newParticipants: Participant[] = [];
    return newParticipants;
};
