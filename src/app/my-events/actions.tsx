'use server';
import type { NewEvent } from '$app/types';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import type { Unpublished } from '@mui/icons-material';
import { doc, getFirestore, setDoc } from 'firebase/firestore';

export const createMyEventDoc = async (docId: string) => {
	const { app, currentUser } = await getAuthorizedAuth();
	if (app && currentUser) {
		const db = getFirestore(app);
		const ref = doc(db, '/users', currentUser.uid, 'my-events', docId);
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
		};
		await setDoc(ref, newEvent);
		return;
	}
	return;
};
