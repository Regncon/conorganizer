// import { Participant } from '$lib/types';
// import { IdTokenResult, User } from 'firebase/auth';
// import { EventTicket } from './actions';
//
// import { processParticipantAssignment } from './AssignParticipant';
//
// describe('processParticipantAssignment', () => {
//     const user: User = {
//         email: 'test@example.com', uid: 'user123',
//         emailVerified: false,
//         isAnonymous: false,
//         metadata: {},
//         providerData: [],
//         refreshToken: '',
//         tenantId: null,
//         delete: function (): Promise<void> {
//             throw new Error('Function not implemented.');
//         },
//         getIdToken: function (forceRefresh?: boolean): Promise<string> {
//             throw new Error('Function not implemented.');
//         },
//         getIdTokenResult: function (forceRefresh?: boolean): Promise<IdTokenResult> {
//             throw new Error('Function not implemented.');
//         },
//         reload: function (): Promise<void> {
//             throw new Error('Function not implemented.');
//         },
//         toJSON: function (): object {
//             throw new Error('Function not implemented.');
//         },
//         displayName: null,
//         phoneNumber: null,
//         photoURL: null,
//         providerId: ''
//     };
//
//     const tickets : EventTicket[] = [
//         {
//             id: 1,
//             order_id: 1001,
//             crm: {
//                 email: 'test@example.com',
//                 first_name: 'John',
//                 last_name: 'Doe',
//                 born: '1990-01-01',
//                 id: 100
//             },
//             category: 'General Admission',
//             category_id: 1,
//         },
//         {
//             id: 2,
//             order_id: 1002,
//             crm: {
//                 email: 'jane@example.com',
//                 first_name: 'Jane',
//                 last_name: 'Smith',
//                 born: '1995-05-15',
//                 id: 101
//             },
//             category: 'VIP',
//             category_id: 2,
//         },
//     ];
//
//     const participants: Partial<Participant>[] = [
//         {
//             id: 'participant1',
//             ticketId: 1,
//             users: ['user456'],
//             ticketEmail: 'test@example.com',
//             orderEmails: ['test@example.com'],
//         },
//     ];
//
//     it('should create new participants and update existing ones', async () => {
//         const result = await processParticipantAssignment(tickets, participants as Participant[], user);
//
//         // New participants
//         expect(result.newParticipants).toHaveLength(1);
//         expect(result.newParticipants[0].ticketId).toBe(2);
//         expect(result.newParticipants[0].users).toContain('user123');
//
//         // Updated participants
//         expect(result.updatedParticipants).toHaveLength(1);
//         expect(result.updatedParticipants[0].users).toContain('user123');
//
//         // User info to be updated
//         expect(result.myUserInfoToBeUpdated.participantIds).toContain('participant1');
//
//         // Ensure no duplicate participant IDs
//         const participantIds = result.myUserInfoToBeUpdated.participantIds;
//         expect(new Set(participantIds).size).toBe(participantIds.length);
//     });
//
//     it('should handle cases with no new tickets', async () => {
//         const updatedTickets = tickets.slice(0, 1); // Only ticket with id 1
//         const result = await processParticipantAssignment(updatedTickets, participants as Participant[], user);
//
//         expect(result.newParticipants).toHaveLength(0);
//         expect(result.updatedParticipants).toHaveLength(1);
//     });
//
//     it('should handle cases with no existing participants', async () => {
//         const result = await processParticipantAssignment(tickets, [], user);
//
//         expect(result.newParticipants).toHaveLength(2);
//         expect(result.updatedParticipants).toHaveLength(0);
//     });
// });
