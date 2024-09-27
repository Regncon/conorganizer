import { Participant } from '$lib/types';
import { User } from 'firebase/auth';
import { AssignUserToParticipant } from './Helpers';

describe('AssignParticipantToUser', () => {
    // Setup
    const partialUser: Partial<User> = {
        email: 'user1@example.com',
        uid: 'user123',
    };
    const user = partialUser as User;

    describe('Given that there are participants that the user should be assigned to', () => {
        describe('When the user is assigned to the participants', () => {
            const expectedNewParticipants: Partial<Participant>[] = [
                {
                    id: '101',
                    users: ['user123'],
                    updatedBy: 'user1@example.com',
                    updateAt: expect.any(String),
                },
                {
                    id: '102',
                    users: ['user456', 'user123'],
                    updatedBy: 'user1@example.com',
                    updateAt: expect.any(String),
                },
            ];

            const existingParticipants: Participant[] = [
                {
                    id: '101',
                    ticketId: 1,
                    ticketEmail: 'user1@example.com',
                    firstName: 'New Participant 1',
                    lastName: 'User email is in order emails and should be assigned to this participant',
                    over18: true,
                    orderId: 1001,
                    orderEmails: ['user1@example.com'],
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: [],
                    users: [],
                    createdAt: '',
                    createdBy: 'user1@example.com',
                    updateAt: '',
                    updatedBy: 'user1@example.com',
                },
                {
                    id: '102',
                    ticketId: 2,
                    ticketEmail: 'user2@example.com',
                    firstName: 'New Participant 2',
                    lastName: 'User email is in connected emails, should be assigned',
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
                // Existing participants that should not be updated
                {
                    id: '103',
                    ticketId: 3,
                    ticketEmail: 'user1@example.com',
                    firstName: 'New Participant 3',
                    lastName: 'User already assigned to this participant, should not be changed',
                    orderEmails: ['user1@example.com'],
                    users: ['user123'],
                    over18: true,
                    orderId: 1003,
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: [],
                    createdAt: '',
                    createdBy: 'user1@example.com',
                    updateAt: '',
                    updatedBy: 'user1@example.com',
                },
                {
                    id: '104',
                    ticketId: 4,
                    ticketEmail: 'user4@example.com',
                    firstName: 'New Participant 4',
                    lastName: 'User email is not related, should not be updated',
                    orderEmails: ['user4@example.com'],
                    users: ['user456'],
                    over18: true,
                    orderId: 1004,
                    ticketCategory: '',
                    ticketCategoryId: 0,
                    connectedEmails: ['user5@example.com'],
                    createdAt: '',
                    createdBy: 'user4@example.com',
                    updateAt: '',
                    updatedBy: 'user4@example.com',
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
