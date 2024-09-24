import { Participant } from '$lib/types';
import { User } from 'firebase/auth';
import { EventTicket } from './actions';
import { NewTickets } from './NewTickets';

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
                createdBy: 'user1@example.com',
                updateAt: '',
                updatedBy: 'user1@example.com',
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
                createdBy: 'user1@example.com',
                updateAt: '',
                updatedBy: 'user1@example.com',
            },
            {
                ticketId: 3,
                ticketEmail: 'user2@example.com',
                firstName: 'New Participant 3',
                lastName: 'Same orderId as ticket 3 should not be included',
                users: ['user123'],
                over18: false,
                orderId: 1002,
                orderEmails: ['user1@example.com', 'user2@example.com'],
                ticketCategory: '',
                ticketCategoryId: 0,
                connectedEmails: [],
                createdAt: '',
                createdBy: 'user1@example.com',
                updateAt: '',
                updatedBy: 'user1@example.com',
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
                        last_name: 'Same orderId as ticket 3 should not be included',
                        id: 0,
                        born: '01-01-2010',
                    },
                    category: '',
                    category_id: 0,
                },
                // {
                //     id: 4,
                //     order_id: 1003,
                //     crm: {
                //         email: 'user2@example.com',
                //         first_name: 'New Participant 4',
                //         last_name: 'Same orderId as ticket 3 should not be included',
                //         id: 0,
                //         born: '01-01-1990',
                //     },
                //     category: '',
                //     category_id: 0,
                // },
                // {
                //     id: 5,
                //     order_id: 1004,
                //     crm: {
                //         email: 'user3@example.com',
                //         first_name: 'New Participant 5',
                //         last_name: 'Not assosiated with email, new ticket should be ignored',
                //         id: 0,
                //         born: '',
                //     },
                //     category: '',
                //     category_id: 0,
                // },
                // {
                //     id: 6,
                //     order_id: 1005,
                //     crm: {
                //         email: 'user4@example.com',
                //         first_name: 'New Participant 6',
                //         last_name:
                //             'No participant connected and not assosiated with email, new ticket should be ignored',
                //         id: 0,
                //         born: '',
                //     },
                //     category: '',
                //     category_id: 0,
                // },
            ];

            const existingParticipants: Participant[] = [
                // {
                //     ticketId: 5,
                //     ticketEmail: 'user3@example.com',
                //     id: '',
                //     firstName: 'Participant 5',
                //     lastName: 'Not assosiated with email, should be ignored',
                //     users: ['user456'],
                //     over18: false,
                //     orderId: 0,
                //     orderEmails: [],
                //     ticketCategory: '',
                //     ticketCategoryId: 0,
                //     connectedEmails: [],
                //     createdAt: '',
                //     createdBy: '',
                //     updateAt: '',
                //     updatedBy: '',
                // },
            ];

            const result = await NewTickets(tickets, existingParticipants, user);

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
});
