'use server';

import { z } from 'zod';
import type { InitialRegisterFormState } from '../components/RegisterPage';
import type { RegisterDetails } from '$lib/firebase/firebase';

export const validateRegisterFormAction = async (formData: FormData): Promise<InitialRegisterFormState> => {
    const { email, password, confirm } = Object.fromEntries(formData) as RegisterDetails;

    const schemaEmail = z.string().email({ message: 'Ugyldig e-post' });
    const schemaPassword = z
        .object({
            password: z.string().min(6, { message: 'Passordet må innehalde minst 6 teikn' }),
            confirm: z.string().min(6, { message: 'Passordet må innehalde minst 6 teikn' }),
        })
        .required()
        .refine(({ password, confirm: confirmPassword }) => password === confirmPassword, {
            message: 'Passord og stadfestingspassord må vere like',
        });
    const passwordObject: z.infer<typeof schemaPassword> = {
        password,
        confirm,
    };
    const resultEmail = schemaEmail.safeParse(email);
    const resultPasswords = schemaPassword.safeParse(passwordObject);

    const emailError = resultEmail.error?.format()._errors[0];

    const passwordError = resultPasswords.error?.format().password?._errors[0];
    const confirmError = resultPasswords.error?.format().confirm?._errors[0];
    const matchError = resultPasswords.error?.format()?._errors[0];

    const resetErrors: InitialRegisterFormState = {
        emailError: '',
        passwordError: '',
        confirmError: '',
    };

    if (!resultEmail.success || !resultPasswords.success) {
        return {
            emailError: emailError ?? resetErrors.emailError,
            passwordError: (passwordError ? passwordError : matchError) ?? resetErrors.passwordError,
            confirmError: (confirmError ? confirmError : matchError) ?? resetErrors.confirmError,
        };
    }
    console.log(resultEmail.success || resultPasswords.success);

    return resetErrors;
};
