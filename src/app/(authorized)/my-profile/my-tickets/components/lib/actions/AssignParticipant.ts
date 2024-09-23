'use server';
import { getMyUserInfo } from '$app/(authorized)/my-events/lib/actions';
import { GetAllParticipants } from '$app/(public)/components/lib/serverAction';
import { adminDb, getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ActionResponse, ConUser, Participant } from '$lib/types';
import { throws } from 'assert';
import { doc, updateDoc } from 'firebase/firestore';
import { revalidatePath } from 'next/cache';
import { EventTicket, generateParticipant, GetParticipantsByEmail, GetTicketsByEmail } from './actions';
import { Firestore } from 'firebase-admin/firestore';
import { User } from 'firebase/auth';

export const AssignParticipantByEmail = async () => {
    const { db, user } = await getAuthorizedAuth();
    if (db === null || user === null || user.email === null) {
        throw new Error('Failed to get authorized auth');
    }

    const tickets = await GetTicketsByEmail(user.email);
    if (!tickets) {
        throw new Error('Failed to get tickets from Checkin');
    }

    const participants = await GetParticipantsByEmail(user.email);

    // Process data
    const { newParticipants, updatedParticipants, myUserInfoToBeUpdated } = await processParticipantAssignment(
        tickets,
        participants,
        user
    );

    // Database operations
    const newParticipantIds: string[] = [];

    // Add new participants
    for (const newParticipant of newParticipants) {
        const docRef = await adminDb.collection('participants').add(newParticipant);
        newParticipantIds.push(docRef.id);
        console.log('Participant written with ID: ', docRef.id);
    }

    // Update existing participants
    for (const participant of updatedParticipants) {
        await adminDb.collection('participants').doc(participant.id).update(participant);
    }

    // Update or set user info
    const myUserInfo = await getMyUserInfo(db, user);
    if (myUserInfo) {
        await adminDb.collection('users').doc(user.uid).update(myUserInfoToBeUpdated);
    } else {
        await adminDb.collection('users').doc(user.uid).set(myUserInfoToBeUpdated);
    }

    console.log('Completed AssignParticipantByEmail');
    return participants;
};

export const processParticipantAssignment = async (tickets: EventTicket[], participants: Participant[], user: User) => {
    const newTickets = tickets.filter(
        (ticket) => !participants.some((participant) => participant.ticketId === ticket.id)
    );

    const newParticipants: Partial<Participant>[] = [];
    const updatedParticipants: Participant[] = [];

    newTickets.forEach(async (newTicket) => {
        let newParticipant = await generateParticipant(newTicket.id, tickets, user.email as string);
        newParticipant = { ...newParticipant, users: [user.uid] };
        newParticipants.push(newParticipant);
    });

    const myUserInfoToBeUpdated = {
        admin: false,
        participantIds: [] as string[],
    };

    participants.forEach((participant) => {
        if (!participant.users?.includes(user.uid)) {
            participant.users = participant.users || [];
            participant.users.push(user.uid);
            updatedParticipants.push(participant);
        }
        if (!myUserInfoToBeUpdated.participantIds.includes(participant.id!)) {
            myUserInfoToBeUpdated.participantIds.push(participant.id!);
        }
    });

    return {
        newParticipants,
        updatedParticipants,
        myUserInfoToBeUpdated,
    };
};
