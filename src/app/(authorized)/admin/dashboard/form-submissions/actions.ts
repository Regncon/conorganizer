'use server';

import { getMyUserInfo } from '$app/(authorized)/my-events/actions';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { collectionGroup, documentId, getDoc, getDocs, query, where, type Firestore } from 'firebase/firestore';
import { redirect } from 'next/navigation';

const setReadStatus = async (db: Firestore, queryValue: string) => {
    // const myEventsQuery = query(collectionGroup(db, 'my-events'));
    const myEventsQuery = query(collectionGroup(db, 'my-events'), where(documentId(), '==', queryValue));
    console.log(myEventsQuery);

    console.log(await getDocs(myEventsQuery));
};
export const updateReadStatus = async (queryValue: string) => {
    const { app, user, auth, db } = await getAuthorizedAuth();
    if (app !== null && user !== null && auth !== null && db !== null) {
        setReadStatus(db, queryValue);
    }
};
