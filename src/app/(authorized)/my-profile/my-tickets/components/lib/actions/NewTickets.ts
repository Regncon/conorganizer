import { Participant } from '$lib/types';
import { User } from 'firebase/auth';
import { EventTicket } from './actions';
import { generateParticipant } from './Helpers';

export const NewTickets = async (tickets: EventTicket[], participants: Participant[], user: User) => {
    console.log('Assigning new tickets to participants');
    // console.log('tickes', tickets, 'participants', participants, 'user', user);

    const ticketsWithUsersEmail = tickets?.filter((ticket) => ticket.crm.email === (user.email as string));
    const ticketsWithOrderNumberFromEmail = tickets?.filter((ticket) =>
        ticketsWithUsersEmail?.some((emailTicket) => emailTicket.order_id === ticket.order_id)
    );

    const newParticipants: Participant[] = [];
    ticketsWithOrderNumberFromEmail.forEach((ticket) => {
        const newParticipant = generateParticipant(ticket.id, tickets, user.email as string);
        newParticipant.users = [user.uid];
        newParticipants.push(newParticipant);
    });
    return newParticipants;
};
