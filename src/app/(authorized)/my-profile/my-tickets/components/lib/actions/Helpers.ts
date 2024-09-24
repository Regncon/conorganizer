import { ConUser, Participant } from '$lib/types';
import { User } from 'firebase/auth';
import { EventTicket } from './actions';
import { Update } from '@mui/icons-material';

export const NewTickets = (tickets: EventTicket[], participants: Participant[], user: User) => {
    console.log('Assigning new tickets to participants');

    const ticketsWithoutParticipants = tickets?.filter(
        (ticket) => !participants.some((participant) => participant.ticketId === ticket.id)
    );

    const usersTickets = GetTicketsThatBelongToUser(ticketsWithoutParticipants, user);

    const newParticipants: Participant[] = [];
    usersTickets.forEach((ticket) => {
        const newParticipant = generateParticipant(ticket.id, tickets, user.email as string);
        newParticipant.users = [user.uid];
        newParticipants.push(newParticipant);
    });
    return newParticipants;
};

export const AssignUserToParticipant = (participants: Participant[], user: User): Partial<Participant>[] => {
    const filteredParticipants = participants.filter(
        (participant) =>
            (participant.orderEmails.includes(user.email as string) ||
                participant.connectedEmails?.includes(user.email as string)) &&
            !participant.users?.includes(user.uid)
    );

    const assignedParticipants = filteredParticipants.map((participant) => {
        const updatedParticipant: Partial<Participant> = {
            id: participant.id,
            ticketId: participant.ticketId,
            users: [...(participant.users || []), user.uid],
            updatedBy: user.email || '',
            updateAt: new Date().toISOString(),
        };
        return updatedParticipant;
    });

    return assignedParticipants;
};

const GetTicketsThatBelongToUser = (tickets: EventTicket[], user: User) => {
    const ticketsWithUsersEmail = tickets?.filter((ticket) => ticket.crm.email === (user.email as string));

    const ticketsWithOrderNumberFromEmail = tickets?.filter((ticket) =>
        ticketsWithUsersEmail?.some((emailTicket) => emailTicket.order_id === ticket.order_id)
    );

    return ticketsWithOrderNumberFromEmail;
};

export const generateParticipant = (ticketId: number, tickets: EventTicket[], userEmail: string) => {
    const ticket = tickets.find((ticket) => ticket.id === ticketId);
    if (!ticket) throw new Error('ticket not found');

    const isOver18 = new Date().getFullYear() - new Date(ticket.crm.born).getFullYear() > 18;

    const orderEmails = tickets.filter((t) => t.order_id === ticket.order_id).map((t) => t.crm.email);

    let participant: Participant = {
        firstName: ticket.crm.first_name,
        lastName: ticket.crm.last_name,
        over18: isOver18,
        ticketEmail: ticket.crm.email,
        orderEmails: orderEmails,
        ticketId: ticket.id,
        orderId: ticket.order_id,
        ticketCategory: ticket.category,
        ticketCategoryId: ticket.category_id,
        createdAt: new Date().toISOString(),
        createdBy: userEmail,
        updateAt: new Date().toISOString(),
        updatedBy: userEmail,
        connectedEmails: [],
    };
    return participant;
};
