import { Participant } from '$lib/types';
import { User } from 'firebase/auth';
import { EventTicket } from './actions';
import { NewTickets } from './Helpers';

describe('NewTickets', () => {
    // Setup
    const partialUser: Partial<User> = {
        email: 'user1@example.com',
        uid: 'user123',
    };
    const user = partialUser as User;

    describe('Given there are tickets with a users email that is not assigned a participants', () => {
        const expectedNewParticipants: Participant[] = [
            {
                ticketId: 1,
                ticketEmail: 'user1@example.com',
                firstName: 'New Participant 1',
                lastName: 'New ticket should be created',
                users: ['user123'],
                over18: true,
                orderId: 1001,
                orderEmails: ['user1@example.com'],
                ticketCategory: '',
                ticketCategoryId: 0,
                connectedEmails: [],
                createdAt: '',
                createdBy: 'user123',
                updateAt: '',
                updatedBy: 'user123',
            },
            {
                ticketId: 2,
                ticketEmail: 'user1@example.com',
                firstName: 'New Participant 2',
                lastName: 'New ticket should be created',
                orderEmails: ['user1@example.com', 'user2@example.com'],
                users: ['user123'],
                over18: true,
                orderId: 1002,
                ticketCategory: '',
                ticketCategoryId: 0,
                connectedEmails: [],
                createdAt: '',
                createdBy: 'user123',
                updateAt: '',
                updatedBy: 'user123',
            },
            {
                ticketId: 3,
                ticketEmail: 'user2@example.com',
                firstName: 'New Participant 3',
                lastName: 'Same orderId as ticket 2 should be creaded',
                users: ['user123'],
                over18: false,
                orderId: 1002,
                orderEmails: ['user1@example.com', 'user2@example.com'],
                ticketCategory: '',
                ticketCategoryId: 0,
                connectedEmails: [],
                createdAt: '',
                createdBy: 'user123',
                updateAt: '',
                updatedBy: 'user123',
            },
        ];

        describe('when we compear the tickets with the participants', async () => {
            // Test
            const tickets: EventTicket[] = [
                {
                    id: 1,
                    order_id: 1001,
                    crm: {
                        email: 'user1@example.com',
                        first_name: 'New Participant 1',
                        last_name: 'New ticket should be created',
                        id: 0,
                        born: '01-01-1990',
                    },
                    category: '',
                    category_id: 0,
                },
                {
                    id: 2,
                    order_id: 1002,
                    crm: {
                        email: 'user1@example.com',
                        first_name: 'New Participant 2',
                        last_name: 'New ticket should be created',
                        id: 0,
                        born: '01-01-1990',
                    },
                    category: '',
                    category_id: 0,
                },
                {
                    id: 3,
                    order_id: 1002,
                    crm: {
                        email: 'user2@example.com',
                        first_name: 'New Participant 3',
                        last_name: 'Same orderId as ticket 2 should be creaded',
                        id: 0,
                        born: '01-01-2010',
                    },
                    category: '',
                    category_id: 0,
                },
                {
                    id: 4,
                    order_id: 1003,
                    crm: {
                        email: 'user2@example.com',
                        first_name: 'New Participant 4',
                        last_name: 'Not the same email or order Should be ignored',
                        id: 0,
                        born: '01-01-1990',
                    },
                    category: '',
                    category_id: 0,
                },
                {
                    id: 5,
                    order_id: 1004,
                    crm: {
                        email: 'user1@example.com',
                        first_name: 'Participant 5',
                        last_name: 'Ticket already claimed by user and has participant. Ticket should be ignored',
                        id: 0,
                        born: '01-01-1990',
                    },
                    category: '',
                    category_id: 0,
                },
                {
                    id: 6,
                    order_id: 1005,
                    crm: {
                        email: 'user4@example.com',
                        first_name: 'New Participant 6',
                        last_name:
                            'No participant connected and not assosiated with email, new ticket should be ignored',
                        id: 0,
                        born: '',
                    },
                    category: '',
                    category_id: 0,
                },
                {
                    id: 7,
                    order_id: 1002,
                    crm: {
                        email: 'user1@example.com',
                        first_name: 'New Participant 7',
                        last_name: 'Type is dinner, should be ignored',
                        id: 0,
                        born: '01-01-1990',
                    },
                    category: 'Middag',
                    category_id: 157059,
                },
            ];

            const existingParticipants: Participant[] = [
                {
                    ticketId: 5,
                    ticketEmail: 'user1@example.com',
                    firstName: 'Participant 5',
                    lastName: 'Ticket already claimed by user and has participant. Ticket should be ignored',
                    users: ['user123'],
                    over18: true,
                    orderId: 1006,
                    orderEmails: ['user1@example.com'],
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: [],
                    createdAt: '',
                    createdBy: '',
                    updateAt: '',
                    updatedBy: '',
                },
                {
                    ticketId: 6,
                    ticketEmail: 'user4@example.com',
                    firstName: 'New Participant 6',
                    lastName: 'No participant connected and not assosiated with email, new ticket should be ignored',
                    users: [],
                    over18: false,
                    orderId: 1005,
                    orderEmails: [],
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: [],
                    createdAt: '',
                    createdBy: '',
                    updateAt: '',
                    updatedBy: '',
                },
            ];

            const result = NewTickets(tickets, existingParticipants, user);

            test('then the new participants should be created and the existing ones should be updated', async () => {
                const resultWithoutDate = result.map((participant) => ({
                    ...participant,
                    createdAt: undefined, // Ignore `createdAt` in comparison
                    updateAt: undefined, // Ignore `updateAt` in comparison
                }));

                const expectedWithoutDate = expectedNewParticipants.map((participant) => ({
                    ...participant,
                    createdAt: undefined, // Ignore `createdAt` in comparison
                    updateAt: undefined, // Ignore `updateAt` in comparison
                }));

                expect(resultWithoutDate).toMatchObject(expectedWithoutDate);
            });
        });
    });
    describe('Given there are no tickets with a users email that is not assigned a participants', () => {
        describe('when we compear the tickets with the participants', async () => {
            const expectedNewParticipants: Participant[] = [];
            const tickets: EventTicket[] = [
                {
                    id: 1,
                    order_id: 1001,
                    crm: {
                        email: 'user2@example.com',
                        first_name: 'Participant 2',
                        last_name: 'New ticket should be ignored',
                        id: 0,
                        born: '01-01-1990',
                    },
                    category: '',
                    category_id: 0,
                },
                {
                    id: 2,
                    order_id: 1002,
                    crm: {
                        email: 'user3@example.com',
                        first_name: 'Participant 3',
                        last_name: 'New ticket should be ignored',
                        id: 0,
                        born: '01-01-1990',
                    },
                    category: '',
                    category_id: 0,
                },
            ];
            const existingParticipants: Participant[] = [
                {
                    ticketId: 1,
                    ticketEmail: 'user2@example.com',
                    firstName: 'Participant 2',
                    lastName: 'Ticket should be ignored',
                    users: [],
                    over18: true,
                    orderId: 1001,
                    orderEmails: ['user2@example.com'],
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: [],
                    createdAt: '',
                    createdBy: '',
                    updateAt: '',
                    updatedBy: '',
                },
                {
                    ticketId: 2,
                    ticketEmail: 'user3@example.com',
                    firstName: 'Participant 3',
                    lastName: 'New ticket should be ignored',
                    orderEmails: ['user3@example.com'],
                    users: [],
                    over18: true,
                    orderId: 1002,
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: [],
                    createdAt: '',
                    createdBy: '',
                    updateAt: '',
                    updatedBy: '',
                },
            ];
            const result = await NewTickets(tickets, existingParticipants, user);
            test('then no new participants should be created and no existing ones should be updated', async () => {
                expect(result).toMatchObject(expectedNewParticipants);
                expect(result).toHaveLength(0);
            });
        });
    });
});
