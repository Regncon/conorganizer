// participants.test.ts

import { Participant } from '$lib/types';
import { describe, it } from 'node:test';
import { processParticipantAssignment } from './AssignParticipant';

describe('processParticipantAssignment', () => {
    const user = { email: 'test@example.com', uid: 'user123' };

    const tickets = [
        {
            id: 1,
            order_id: 'order1',
            crm: {
                email: 'test@example.com',
                first_name: 'John',
                last_name: 'Doe',
                born: '1990-01-01',
            },
            category: 'General Admission',
            category_id: 'cat1',
        },
        {
            id: 2,
            order_id: 'order2',
            crm: {
                email: 'jane@example.com',
                first_name: 'Jane',
                last_name: 'Smith',
                born: '1995-05-15',
            },
            category: 'VIP',
            category_id: 'cat2',
        },
    ];

    const participants: Partial<Participant>[] = [
        {
            id: 'participant1',
            ticketId: 1,
            users: ['user456'],
            ticketEmail: 'test@example.com',
            orderEmails: ['test@example.com'],
        },
    ];

    it('should create new participants and update existing ones', () => {
        const result = processParticipantAssignment(tickets, participants, user);

        // New participants
        expect(result.newParticipants).toHaveLength(1);
        expect(result.newParticipants[0].ticketId).toBe(2);
        expect(result.newParticipants[0].users).toContain('user123');

        // Updated participants
        expect(result.updatedParticipants).toHaveLength(1);
        expect(result.updatedParticipants[0].users).toContain('user123');

        // User info to be updated
        expect(result.myUserInfoToBeUpdated.participantIds).toContain('participant1');

        // Ensure no duplicate participant IDs
        const participantIds = result.myUserInfoToBeUpdated.participantIds;
        expect(new Set(participantIds).size).toBe(participantIds.length);
    });

    it('should handle cases with no new tickets', () => {
        const updatedTickets = tickets.slice(0, 1); // Only ticket with id 1
        const result = processParticipantAssignment(updatedTickets, participants, user);

        expect(result.newParticipants).toHaveLength(0);
        expect(result.updatedParticipants).toHaveLength(1);
    });

    it('should handle cases with no existing participants', () => {
        const result = processParticipantAssignment(tickets, [], user);

        expect(result.newParticipants).toHaveLength(2);
        expect(result.updatedParticipants).toHaveLength(0);
    });
});
