import { Participant } from '$lib/types';
import { User } from 'firebase/auth';
import { AssignUserToParticipant } from './NewTickets';

describe('AssignParticipantToUser', () => {
    // Setup
    const partialUser: Partial<User> = {
        email: 'user1@example.com',
        uid: 'user123',
    };
    const user = partialUser as User;
    describe('Given that there are pargticipants that the user should be assigned to', () => {
        describe('When the user is assigned to the participants', () => {
            const expectedNewParticipants: Participant[] = [
                {
                    ticketId: 1,
                    ticketEmail: 'user1@example.com',
                    firstName: 'New Participant 1',
                    lastName: 'User shold be assigned to this participant',
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
                // {
                //     ticketId: 2,
                //     ticketEmail: 'user1@example.com',
                //     firstName: 'New Participant 2',
                //     lastName: 'New ticket should be created',
                //     orderEmails: ['user1@example.com', 'user2@example.com'],
                //     users: ['user123'],
                //     over18: true,
                //     orderId: 1002,
                //     ticketCategory: '',
                //     ticketCategoryId: 0,
                //     connectedEmails: [],
                //     createdAt: '',
                //     createdBy: 'user1@example.com',
                //     updateAt: '',
                //     updatedBy: 'user1@example.com',
                // },
                // {
                //     ticketId: 3,
                //     ticketEmail: 'user2@example.com',
                //     firstName: 'New Participant 3',
                //     lastName: 'Same orderId as ticket 2 should be creaded',
                //     users: ['user123'],
                //     over18: false,
                //     orderId: 1002,
                //     orderEmails: ['user1@example.com', 'user2@example.com'],
                //     ticketCategory: '',
                //     ticketCategoryId: 0,
                //     connectedEmails: [],
                //     createdAt: '',
                //     createdBy: 'user1@example.com',
                //     updateAt: '',
                //     updatedBy: 'user1@example.com',
                // },
            ];
            const existingParticipants: Participant[] = [
                {
                    ticketId: 1,
                    ticketEmail: 'user1@example.com',
                    firstName: 'New Participant 1',
                    lastName: 'User shold be assigned to this participant',
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
            ];
            const result = AssignUserToParticipant(existingParticipants, user);

            test('Then the user should be assigned to the participants', () => {
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
