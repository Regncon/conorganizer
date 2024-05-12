'use client';

import { logout, setSessionCookie } from '$app/(auth)/login/action';

import { initializeApp } from 'firebase/app';
import {
    getAuth,
    signInWithEmailAndPassword,
    signOut,
    createUserWithEmailAndPassword,
    sendPasswordResetEmail,
} from 'firebase/auth';
import { getFirestore } from 'firebase/firestore';
import type { FormEvent } from 'react';
import { firebaseConfig } from './config';

const app = initializeApp(firebaseConfig, 'client');

export const db = getFirestore(app);
export const firebaseAuth = getAuth(app);

type LoginDetails = {
    email: string;
    password: string;
};
export type RegisterDetails = {
    email: string;
    password: string;
    confirm: string;
};

const getEmailAndPasswordFromFormData = (e: FormEvent<HTMLFormElement>) => {
    const target = e.target as HTMLFormElement;
    const { email, password } = Object.fromEntries(new FormData(target)) as LoginDetails;
    return { email, password };
};
export const signInAndCreateCookie = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const { email, password } = getEmailAndPasswordFromFormData(e);
    if (!!!email && !!!password) {
        return;
    }
    const userCredentials = await signInWithEmailAndPassword(firebaseAuth, email, password);
    const idToken = await userCredentials.user.getIdToken();

    await setSessionCookie(idToken);
};
export const signOutAndDeleteCookie = async () => {
    await signOut(firebaseAuth);
    await logout();
};

export const singUpAndCreateCookie = async (e: FormEvent<HTMLFormElement>) => {
    const { email, password } = getEmailAndPasswordFromFormData(e);
    const userCredentials = await createUserWithEmailAndPassword(firebaseAuth, email, password);
    const idToken = await userCredentials.user.getIdToken();

    await setSessionCookie(idToken);
};
export const forgotPassword = async (email: string) => {
    const forgot = await sendPasswordResetEmail(firebaseAuth, email);
};

// export const eventsRef = collection(db, 'events');
// export const eventRef = (id: string) => doc(db, `events/${id}`);
// // export const allUserSettingsRef = collection(db, 'usersettings');
// // export const userSettingsRef = (userId: string) => doc(db, `usersettings/${userId}`);
// // export const participantEnrollmentsRef = (eventId: string, userId: string, participantId: string ) =>
// //     doc(db, `events/${eventId}`, `/enrollments/${userId}`,`/eventParticipants/${participantId}`);
// // export const participantsRef = (userId: string ) =>  collection(db, `usersettings/${userId}/participants/`);
// // export const participantRef = (userId: string, participantId: string ) => doc(db, `usersettings/${userId}/participants/${participantId}`);
// // export const enrollmentChoicesRef = (eventId: string ) => collection(db, `events/${eventId}/enrollmentChoices/`);
