import { collectionData, docData } from 'rxfire/firestore';
import { authorizationRef, eventRef, eventsRef } from './firebase';

export const allEvents$ = collectionData(eventsRef, { idField: 'id' });
export function singleEvent$(id: string) {
    return docData(eventRef(id), { idField: 'id' });
}
export function conAuthorizationRef$(userId: string) {
    return docData(authorizationRef(userId), { idField: 'id' });
}
