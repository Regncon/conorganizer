import { collection } from 'firebase/firestore';
import { collectionData } from 'rxfire/firestore';
import { db } from '../firebase';

export const eventsRef = collection(db, 'events');
export const allEvents$ = collectionData(eventsRef, { idField: 'id' });
