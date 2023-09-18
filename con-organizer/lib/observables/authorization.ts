import { doc } from 'firebase/firestore';
import { doc as rxDoc } from 'rxfire/firestore';
import { db } from '../firebase';

export const authorizationRef = (userId: string) => doc(db, `authorization/${userId}`);
export function conAuthorizationRef$(userId: string) {
    return rxDoc(authorizationRef(userId));
}
