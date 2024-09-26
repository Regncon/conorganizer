'use client';

import { logout, setSessionCookie } from '$app/(auth)/login/lib/actions';
import { initializeApp } from 'firebase/app';
import {
    getAuth,
    signInWithEmailAndPassword,
    signOut,
    createUserWithEmailAndPassword,
    sendPasswordResetEmail,
} from 'firebase/auth';
import { getFirestore } from 'firebase/firestore';
import { firebaseConfig } from './config';

const app = initializeApp(firebaseConfig, 'client');

export const db = getFirestore(app);
export const firebaseAuth = getAuth(app);

export type LoginDetails = {
    email: string;
    password: string;
};
export type RegisterDetails = {
    email: string;
    password: string;
    confirm: string;
};

const getEmailAndPasswordFromFormData: (formData: FormData) => LoginDetails = (formData) => {
    const { email, password } = Object.fromEntries(formData) as LoginDetails;
    return { email, password };
};
export const signInAndCreateCookie: (formData: FormData) => Promise<void> = async (formData) => {
    const { email, password } = getEmailAndPasswordFromFormData(formData);
    if (!!!email && !!!password) {
        return;
    }
    const userCredentials = await signInWithEmailAndPassword(firebaseAuth, email, password);
    const idToken = await userCredentials.user.getIdToken();

    await setSessionCookie(idToken);
};

export const signOutAndDeleteCookie: () => Promise<void> = async () => {
    clearLocalStorage();
    await signOut(firebaseAuth);
    await logout();
};

const clearLocalStorage: () => void = () => {
    document.cookie = 'myParticipants=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    localStorage.removeItem('filters');
};

export const singUpAndCreateCookie: (formData: FormData) => Promise<void> = async (formData) => {
    const { email, password } = getEmailAndPasswordFromFormData(formData);
    const userCredentials = await createUserWithEmailAndPassword(firebaseAuth, email, password);
    const idToken = await userCredentials.user.getIdToken();

    await setSessionCookie(idToken);
};
export const forgotPassword: (email: string) => Promise<void> = async (email) => {
    await sendPasswordResetEmail(firebaseAuth, email);
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
