'use server';

import { SESSION_COOKIE_NAME, adminAuth } from '$lib/firebase/firebaseAdmin';
import type { ResponseCookie } from 'next/dist/compiled/@edge-runtime/cookies';
import { cookies } from 'next/headers';
import { z } from 'zod';
import type { InitialLoginFormState } from '../components/LoginPage';

export const setSessionCookie = async (idToken: string) => {
    console.log(idToken, 'idToken');
    const cookieStore = cookies();
    const twoWeekExpire = 14 * 24 * 60 * 60 * 1000;
    const expirationDate = Date.now() + twoWeekExpire;
    const twoWeeksInSeconds = 14 * 24 * 60 * 60;

    const adminIdToken = await adminAuth.createSessionCookie(idToken, { expiresIn: twoWeekExpire });
    console.log('adminIdToken', adminIdToken);
    const options: Partial<ResponseCookie> = {
        maxAge: twoWeeksInSeconds,
        httpOnly: true,
        secure: true,
        expires: expirationDate,
    };
    console.log(new Date(expirationDate), expirationDate, options);

    cookieStore.set(SESSION_COOKIE_NAME, adminIdToken, options);
};

export const logout = async () => {
    const cookieStore = cookies();
    cookieStore.delete(SESSION_COOKIE_NAME);
};

export const validateLoginForm = async (formData: FormData): Promise<InitialLoginFormState> => {
    const formDataEntries = Object.fromEntries(formData);

    const schemaEmail = z.string().email({ message: 'Ugyldig e-post' });
    const schemaPassword = z.string().min(6, { message: 'Passordet må innehalde minst 6 teikn' });

    const resultEmail = schemaEmail.safeParse(formDataEntries.email);
    const resultPassword = schemaPassword.safeParse(formDataEntries.password);

    const resetErrors: InitialLoginFormState = {
        emailError: '',
        passwordError: '',
    };

    if (!resultEmail.success || !resultPassword.success) {
        return {
            passwordError: resultPassword.error?.issues[0].message ?? resetErrors.passwordError,
            emailError: resultEmail.error?.issues[0].message ?? resetErrors.emailError,
        };
    }

    return resetErrors;
};
