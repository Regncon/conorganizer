import { doc } from 'firebase/firestore';
import { doc as rxDoc } from 'rxfire/firestore';
import { db } from '../firebase';

export const eventRef = (id: string) => doc(db, `events/${id}`);
export function singleEvent$(id: string) {
    return rxDoc(eventRef(id));
}
