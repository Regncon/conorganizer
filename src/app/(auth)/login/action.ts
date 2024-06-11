'use server';

import { SESSION_COOKIE_NAME, adminAuth } from '$lib/firebase/firebaseAdmin';
import type { ResponseCookie } from 'next/dist/compiled/@edge-runtime/cookies';
import { cookies } from 'next/headers';

export const setSessionCookie = async (idToken: string) => {
    const cookieStore = cookies();
    const twoWeekExpire = 14 * 24 * 60 * 60 * 1000;
    const expirationDate = Date.now() + twoWeekExpire;
    const twoWeeksInSeconds = 14 * 24 * 60 * 60;

    const adminIdToken = await adminAuth.createSessionCookie(idToken, { expiresIn: twoWeekExpire });
    const options: Partial<ResponseCookie> = {
        maxAge: twoWeeksInSeconds,
        httpOnly: true,
        secure: true,
        expires: expirationDate,
    };
    console.log(new Date(expirationDate), expirationDate, options);

    cookieStore.set(SESSION_COOKIE_NAME, adminIdToken, options);
};

export async function logout() {
    const cookieStore = cookies();
    cookieStore.delete(SESSION_COOKIE_NAME);
}
