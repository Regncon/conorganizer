import { DocumentData } from 'firebase/firestore';
import { collectionData, doc, docData } from 'rxfire/firestore';
import { BehaviorSubject, Observable } from 'rxjs';
import { ConAuthorization } from '@/models/types';
import { authorizationRef, eventRef, eventsRef } from './firebase';

export const allEvents$ = collectionData(eventsRef, { idField: 'id' });
export function singleEvent$(id: string) {
    return docData(eventRef(id), { idField: 'id' });
}

const userSettingsObs$ = new BehaviorSubject<Observable<DocumentData | undefined> | undefined>(undefined);
export function conAuthorizationRef$(userId: string) {
    return docData(authorizationRef(userId), { idField: 'id' });
}
// export function conAuthorizationRef$(userId?: string) {
//     if (userId) {
//         userSettingsObs$.next(docData(authorizationRef(userId), { idField: 'id' }));
//         return userSettingsObs$;
//     }
//     return userSettingsObs$;
// }
