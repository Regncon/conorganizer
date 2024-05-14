'use server';
import type { NewEvent } from '$app/types';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import type { Unpublished } from '@mui/icons-material';
import { collection, doc, getDocs, getFirestore, setDoc, type Firestore } from 'firebase/firestore';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import { revalidatePath } from 'next/cache';
import { createUserWithEmailAndPassword, signInWithEmailAndPassword, type User } from 'firebase/auth';
import { firebaseAuth } from '$lib/firebase/firebase';

export const createMyEventDoc = async (docId: string) => {
    const { app, user } = await getAuthorizedAuth();
    if (app && user && user.email) {
        const db = getFirestore(app);
        const ref = doc(db, '/users', user.uid, 'my-events', docId);
        const newEvent: Omit<NewEvent, 'id'> = {
            fridayEvening: true,
            saturdayEvening: true,
            saturdayMorning: true,
            sundayMorning: true,
            unpublished: true,
            additionalComments: '',
            adultsOnly: false,
            beginnerFriendly: false,
            childFriendly: false,
            description: '',
            email: '',
            gameType: '',
            lessThanThreeHours: false,
            moduleCompetition: false,
            moreThanSixHours: false,
            name: '',
            participants: 0,
            phone: '',
            possiblyEnglish: false,
            system: '',
            title: '',
            volunteersPossible: false,
            createdAt: new Date(Date.now()).toString(),
            createdBy: user.email,
            updateAt: '',
            updatedBy: '',
            subTitle: '',
        };
        await setDoc(ref, newEvent);
        return;
    }
    return;
};
export async function getAllMyEvents(db: Firestore, user: User) {
    const ref = collection(db, '/users', user.uid, 'my-events');
    const docs = await getDocs(ref);
    const myEvents = docs.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as NewEvent[];
    return myEvents;
}

export async function updateMyEvents() {
    revalidatePath('/my-events', 'page');
}
