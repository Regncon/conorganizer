'use client';

import { login, logout } from '$app/(auth)/login/action';

import { initializeApp, type FirebaseOptions } from 'firebase/app';
import {
	browserLocalPersistence,
	getAuth,
	setPersistence,
	signInWithEmailAndPassword,
	signOut,
	type Config,
} from 'firebase/auth';

import { collection, doc, getFirestore } from 'firebase/firestore';
import type { FormEvent } from 'react';
import { firebaseConfig } from './config';

const app = initializeApp(firebaseConfig, 'client');

export const db = getFirestore(app);
export const firebaseAuth = getAuth(app);
setPersistence(firebaseAuth, browserLocalPersistence);

type LoginDetails = {
	email: string;
	password: string;
};
export const signInAndCreateCookie = async (e: FormEvent<HTMLFormElement>) => {
	e.preventDefault();
	const target = e.target as HTMLFormElement;
	const { email, password } = Object.fromEntries(new FormData(target)) as LoginDetails;
	console.log();
	if (!!!email && !!!password) {
		return;
	}
	const userCredentials = await signInWithEmailAndPassword(firebaseAuth, email, password);
	const idToken = await userCredentials.user.getIdToken(true);

	login(idToken);
};
export const signOutAndDeleteCookie = () => {
	signOut(firebaseAuth);
	logout();
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
