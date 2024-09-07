'use server';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import { revalidatePath } from 'next/cache';
import type { ConEvent } from '$lib/types';

export async function getAllEvents() {
    const eventRef = await adminDb.collection('events').get();
    const events = eventRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as ConEvent[];
    return events;
}

export async function getEventById(id: string) {
    const event = (await adminDb.collection('events').doc(id).get()).data() as ConEvent;
    return { ...event, id };
}

export async function updateEvents() {
    revalidatePath('/', 'page');
}
export async function updateDashboardEvents() {
    revalidatePath('/admin/dashboard/events', 'page');
}
