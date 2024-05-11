'use server';

import { SESSION_COOKIE_NAME } from '$lib/firebase/firebaseAdmin';
import type { UserCredential } from 'firebase/auth';
import { cookies } from 'next/headers';

export async function login(idToken: string) {
	const cookieStore = cookies();
	// const expiresAfterAYear = new Date();
	// expiresAfterAYear.setTime(expiresAfterAYear.getTime() + 350 * 24 * 60 * 60 * 1000);
	const expiresAfter5Days = new Date();
	expiresAfter5Days.setTime(expiresAfter5Days.getTime() + 60 * 60 * 24 * 5 * 1000);

	cookieStore.set(SESSION_COOKIE_NAME, idToken, {
		expires: expiresAfter5Days,
	});
}
export async function logout() {
	const cookieStore = cookies();
	cookieStore.delete(SESSION_COOKIE_NAME);
}
