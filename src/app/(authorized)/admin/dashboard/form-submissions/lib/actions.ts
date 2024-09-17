'use server';

import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import type { MyNewEvent } from '$lib/types';
import { doc, updateDoc, type DocumentReference, type Firestore } from 'firebase/firestore';

export type MyEventUpdateValueName = { isRead: boolean };
const setReadStatus = async (db: Firestore, queryValue: string, updatedValue: MyEventUpdateValueName) => {
    const userRef = doc(db, queryValue) as DocumentReference<MyNewEvent, MyNewEvent>;
    await updateDoc(userRef, updatedValue);
};
export const updateReadAndOrAcceptedStatus = async (queryValue: string, updatedValue: MyEventUpdateValueName) => {
    const { app, user, auth, db } = await getAuthorizedAuth();
    if (app !== null && user !== null && auth !== null && db !== null) {
        setReadStatus(db, queryValue, updatedValue);
    }
};
