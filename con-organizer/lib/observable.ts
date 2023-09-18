import { collectionData, docData } from 'rxfire/firestore';
import { authorizationRef as userSettingsRef, eventRef, eventsRef } from './firebase';

export const allEvents$ = collectionData(eventsRef, { idField: 'id' });
export function singleEvent$(id: string) {
    return docData(eventRef(id), { idField: 'id' });
}
export function userSettings$(userId: string) {
    return docData(userSettingsRef(userId), { idField: 'id' });
}
