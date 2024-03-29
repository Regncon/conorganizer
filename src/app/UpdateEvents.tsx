'use client';
import { db } from '$lib/firebase/firebase';
import { collection, onSnapshot } from 'firebase/firestore';
import { useEffect } from 'react';
import { updateEvents } from './serverAction';

const UpdateEvents = () => {
	useEffect(() => {
		const eventsRef = collection(db, 'event');
		const unsubscribe = onSnapshot(eventsRef, (snapshot) => {
			updateEvents();
		});

		return () => {
			unsubscribe();
		};
	}, []);
	return null;
};

export default UpdateEvents;
