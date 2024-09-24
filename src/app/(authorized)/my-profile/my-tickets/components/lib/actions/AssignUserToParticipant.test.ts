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
                    lastName: 'User email is in order emails and shold be assigned to this participant',
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
                    ticketEmail: 'user2@example.com',
                    firstName: 'New Participant 2',
                    lastName: 'user email is in connected emails, should be assigned',
                    orderEmails: ['user2@example.com'],
                    users: ['user456', 'user123'],
                    over18: true,
                    orderId: 1002,
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: ['user1@example.com'],
                    createdAt: '',
                    createdBy: 'user2@example.com',
                    updateAt: '',
                    updatedBy: 'user2@example.com',
                },
            ];
            const existingParticipants: Participant[] = [
                {
                    ticketId: 1,
                    ticketEmail: 'user1@example.com',
                    firstName: 'New Participant 1',
                    lastName: 'User email is in order emails and shold be assigned to this participant',
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
                    ticketEmail: 'user2@example.com',
                    firstName: 'New Participant 2',
                    lastName: 'user email is in connected emails, should be assigned',
                    orderEmails: ['user2@example.com'],
                    users: ['user456'],
                    over18: true,
                    orderId: 1002,
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: ['user1@example.com'],
                    createdAt: '',
                    createdBy: 'user2@example.com',
                    updateAt: '',
                    updatedBy: 'user2@example.com',
                },
                {
                    ticketId: 3,
                    ticketEmail: 'user2@example.com',
                    firstName: 'New Participant 3',
                    lastName: 'user email is in not in order emails or connected emails, should not be assigned',
                    orderEmails: ['user2@example.com'],
                    users: ['user456'],
                    over18: true,
                    orderId: 1002,
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: ['user3@example.com'],
                    createdAt: '',
                    createdBy: 'user1@example.com',
                    updateAt: '',
                    updatedBy: 'user1@example.com',
                },
            ];
            console.log('existingParticipants', existingParticipants);
            const result = AssignUserToParticipant(existingParticipants, user);
            console.log('result', result);

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
