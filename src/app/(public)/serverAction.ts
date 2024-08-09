'use server';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import { revalidatePath } from 'next/cache';
import type { ConEvent } from '../../lib/types';
import { createUserWithEmailAndPassword, signInWithEmailAndPassword } from 'firebase/auth';
import { firebaseAuth } from '$lib/firebase/firebase';

export async function getAllEvents() {
    const eventRef = await adminDb.collection('events').get();
    const events = eventRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as ConEvent[];
    return events;
}

export async function updateEvents() {
    revalidatePath('/', 'page');
}
