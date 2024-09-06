'use client';
import EmailTextField from '../../shared/EmailTextField';
import { useRouter, useSearchParams } from 'next/navigation';
import { forgotPassword, type LoginDetails } from '$lib/firebase/firebase';
import { useState } from 'react';
import { Grid2, Typography } from '@mui/material';
import { useFormState } from 'react-dom';
import ForgotPasswordButton from './ui/ForgotPasswordButton';
import { validateForgotFormAction } from '../lib/actions';

export const initialForgotFormState = {
    emailError: '',
};
export type InitialForgotFormState = typeof initialForgotFormState;

const ForgotPasswordPage = () => {
    const [message, setMessage] = useState<string>('');
    const router = useRouter();
    const searchParams = useSearchParams();
    const searchParamEmail = searchParams.get('email') ?? undefined;
    const handleFormSubmit: (formData: FormData) => Promise<void> = async (formData) => {
        const { email } = Object.fromEntries(formData) as Omit<LoginDetails, 'password'>;
        setMessage(
            'Ein lenkje for å tilbakestille passordet er sendt viss du har ei registrert e-postadresse hos oss. Du vil nå bli omdirigert til innloggingssida.'
        );
        await forgotPassword(email);
        await new Promise((resolve) => {
            setTimeout(async () => {
                resolve(router.replace('/login'));
            }, 4000);
        });
    };

    const validateForgotForm: (
        state: InitialForgotFormState,
        formData: FormData
    ) => Promise<InitialForgotFormState> = async (_, formData) => {
        const validatedState = await validateForgotFormAction(formData);
        if (!validatedState.emailError) {
            await handleFormSubmit(formData);
        }
        return validatedState;
    };

    const [state, formAction] = useFormState(validateForgotForm, initialForgotFormState);

    return (
        <>
            <Typography sx={{ textAlign: 'center' }}>{message}</Typography>
            <Grid2
                container
                component="form"
                sx={{
                    placeContent: 'center',
                    placeItems: 'center',
                    minWidth: '20rem',
                    marginBlockStart: '1rem',
                    flexDirection: 'column',
                    gap: '1rem',
                }}
                action={formAction}
            >
                <EmailTextField
                    defaultValue={searchParamEmail ?? undefined}
                    error={!!state.emailError}
                    helperText={state.emailError}
                />

                <ForgotPasswordButton />
            </Grid2>
        </>
    );
};
export default ForgotPasswordPage;
