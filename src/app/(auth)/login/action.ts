'use server';

import { SESSION_COOKIE_NAME, getAuthorizedAuth, testAuth } from '$lib/firebase/firebaseAdmin';
import type { UserCredential } from 'firebase/auth';
import type { ResponseCookie } from 'next/dist/compiled/@edge-runtime/cookies';
import { cookies } from 'next/headers';

export const setSessionCookie = async (idToken: string) => {
	const cookieStore = cookies();
	const twoWeekExpire = 14 * 24 * 60 * 60 * 1000;
	const expirationDate = new Date();
	expirationDate.setTime(expirationDate.getTime() + twoWeekExpire);

	const sessionCookie = await testAuth.createSessionCookie(idToken, { expiresIn: twoWeekExpire });
	const option: Partial<ResponseCookie> = {
		maxAge: twoWeekExpire,
		httpOnly: true,
		secure: true,
		expires: expirationDate,
	};

	cookieStore.set(SESSION_COOKIE_NAME, sessionCookie, option);
};

export async function logout() {
	const cookieStore = cookies();
	cookieStore.delete(SESSION_COOKIE_NAME);
}
