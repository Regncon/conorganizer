'use server';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import { revalidatePath } from 'next/cache';
import type { Event } from './types';
export async function getByID(id: string) {
	const eventRef = adminDb.collection('event').doc(id);
	const doc = await eventRef.get();
	if (!doc.exists) {
		console.log('No such document!');
		return null;
	} else {
		console.log('Document data:', doc.data());
		return doc.data();
	}
}
export async function getAll() {
	const eventRef = await adminDb.collection('event').get();
	const events = eventRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as Event[];
	return events;
}

export async function updateEvents() {
	revalidatePath('/', 'page');
}
