import { Participant } from '$lib/types';
import { EventTicket } from './actions';

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
