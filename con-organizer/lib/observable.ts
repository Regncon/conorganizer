import { collectionData, docData } from 'rxfire/firestore';
import {
    allUserSettingsRef,
    enrollmentChoicesRef,
    eventRef,
    eventsRef,
    participantEnrollmentsRef,
    participantRef,
    participantsRef,
    userSettingsRef,
} from './firebase';

export const allEvents$ = collectionData(eventsRef, { idField: 'id' });
export function singleEvent$(id: string) {
    return docData(eventRef(id), { idField: 'id' });
}

export const allUserSettings$ = collectionData(allUserSettingsRef, { idField: 'id' });
export function userSettings$(userId: string) {
    return docData(userSettingsRef(userId), { idField: 'id' });
}

export function participantEnrollments$(eventId: string, userId: string, participantId: string) {
    return docData(participantEnrollmentsRef(eventId, userId, participantId), { idField: 'id' });
}

export function allParticipants$(userId: string) {
    return collectionData(participantsRef(userId), { idField: 'id' });
}

export function singleParticipant$(userId: string, participantId: string) {
    return docData(participantRef(userId, participantId), { idField: 'id' });
}

export function allEnrollmentChoices$(eventId: string) {
    return collectionData(enrollmentChoicesRef(eventId), { idField: 'id' });
}
