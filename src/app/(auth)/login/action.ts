'use server';

import { SESSION_COOKIE_NAME, adminAuth } from '$lib/firebase/firebaseAdmin';
import type { ResponseCookie } from 'next/dist/compiled/@edge-runtime/cookies';
import { cookies } from 'next/headers';
import { z } from 'zod';
import type { InitialFormState } from './LoginPage';

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

export const logout = async () => {
    const cookieStore = cookies();
    cookieStore.delete(SESSION_COOKIE_NAME);
};

export const validateForm = async (formData: FormData): Promise<InitialFormState> => {
    const formDataEntries = Object.fromEntries(formData);

    const schemaEmail = z.string().email({ message: 'Ugyldig e-post' });
    const schemaPassword = z.string().min(6, { message: 'Passordet må innehalde minst 6 teikn' });

    const resultEmail = schemaEmail.safeParse(formDataEntries.email);
    const resultPassword = schemaPassword.safeParse(formDataEntries.password);

    const resetErrorsAfterSuccessfulValidation: InitialFormState = {
        emailError: '',
        passwordError: '',
    };

    if (!resultEmail.success || !resultPassword.success) {
        return {
            passwordError:
                resultPassword.error?.issues[0].message ?? resetErrorsAfterSuccessfulValidation.passwordError,
            emailError: resultEmail.error?.issues[0].message ?? resetErrorsAfterSuccessfulValidation.emailError,
        };
    }

    return resetErrorsAfterSuccessfulValidation;
};
