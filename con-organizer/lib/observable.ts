import { collectionData, docData } from 'rxfire/firestore';
import { eventRef, eventsRef, userEnrollmentsRef,userSettingsRef as userSettingsRef } from './firebase';

export const allEvents$ = collectionData(eventsRef, { idField: 'id' });
export function singleEvent$(id: string) {
    return docData(eventRef(id), { idField: 'id' });
}
export function userSettings$(userId: string) {
    return docData(userSettingsRef(userId), { idField: 'id' });
}

export function userEnrollments$(eventId: string, userId: string) {
    return docData(userEnrollmentsRef(eventId, userId), { idField: 'id' });
}