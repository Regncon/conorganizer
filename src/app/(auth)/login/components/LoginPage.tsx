'use client';
import { Button, CircularProgress, Grid2, Typography } from '@mui/material';
import PasswordTextField from './ui/PasswordTextField';
import { firebaseAuth, signInAndCreateCookie } from '$lib/firebase/firebase';
import { useEffect, useState, useTransition, type ComponentProps } from 'react';
import EmailTextField from '../../shared/EmailTextField';
import { useRouter, useSearchParams } from 'next/navigation';
import { useFormState } from 'react-dom';
import { validateLoginForm as validateLoginFormAction } from '../lib/actions';
import LoginButton from './ui/LoginButton';
import { updateSearchParamsWithEmail } from '../../shared/utils';
import { getRedirectResult, GoogleAuthProvider, onAuthStateChanged, signInWithRedirect, User } from 'firebase/auth';
import GoogleSignInButton from './ui/GoogleButton';
import { GetMyParticipants } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import { cookies } from 'next/headers';

const initialLoginFormState = {
    emailError: '',
    passwordError: '',
};
export type InitialLoginFormState = typeof initialLoginFormState;
export const disableAndLoadingSpinner = (shouldSpin: boolean, isPending: boolean) => ({
    disabled: isPending,
    endIcon: isPending && shouldSpin ? <CircularProgress size="1.5rem" /> : undefined,
});

const initialSpinnersState = {
    forgot: false,
    register: false,
};

const setMyParticipants = async () => {
    const myParticipants = await GetMyParticipants();
    document.cookie = `myParticipants=${JSON.stringify(myParticipants)}; path=/;`;
};

const LoginPage = () => {
    const [isPending, startTransition] = useTransition();
    const [spinners, setSpinners] = useState<typeof initialSpinnersState>(initialSpinnersState);

    const router = useRouter();
    const searchParams = useSearchParams();
    const email = searchParams.get('email') ?? '';
    const expiredSession = searchParams.get('expired') === 'true' ? true : false;

    useEffect(() => {
        router.prefetch('/');
    }, []);

    const handleFormSubmit = async (formData: FormData) => {
        startTransition(async () => {
            await signInAndCreateCookie(formData);
            await setMyParticipants();
            router.replace('/');
        });
    };

    const handleForgotPasswordClick: ComponentProps<'button'>['onClick'] = async () => {
        setSpinners({ ...spinners, forgot: true });
        startTransition(() => {
            router.push(`/forgot-password?email=${email}`);
        });
    };

    const handleRegisterNewUser = async () => {
        setSpinners({ ...spinners, register: true });
        startTransition(() => {
            router.push(`/register`);
        });
    };

    const validateLoginForm: (
        state: InitialLoginFormState,
        formData: FormData
    ) => Promise<InitialLoginFormState> = async (_, formData) => {
        const validatedState = await validateLoginFormAction(formData);
        if (!validatedState.emailError && !validatedState.passwordError) {
            await handleFormSubmit(formData);
        }
        return validatedState;
    };

    const [state, formAction] = useFormState(validateLoginForm, initialLoginFormState);

    return (
        <>
            {expiredSession && (
                <Typography variant="h1" textAlign="center">
                    Økta har gått ut, ver venleg og logg inn igjen.
                </Typography>
            )}

            <Grid2
                component="form"
                container
                sx={{
                    placeContent: 'center',
                    flexDirection: 'column',
                    minWidth: '20rem',
                    gap: '1rem',
                    paddingBlock: '1rem',
                }}
                onChange={(e) => {
                    updateSearchParamsWithEmail(e, router, '/login');
                }}
                action={formAction}
            >
                <GoogleSignInButton />

                <EmailTextField defaultValue={email} error={!!state.emailError} helperText={state.emailError} />
                <PasswordTextField error={!!state.passwordError} helperText={state.passwordError} />

                <LoginButton disabled={isPending} />

                <Button onClick={handleForgotPasswordClick} {...disableAndLoadingSpinner(spinners.forgot, isPending)}>
                    Gløymd passord?
                </Button>
                <Button
                    fullWidth
                    sx={{ marginLeft: 'auto', marginRight: 'auto' }}
                    onClick={handleRegisterNewUser}
                    {...disableAndLoadingSpinner(spinners.register, isPending)}
                >
                    Registrer ny brukar
                </Button>
            </Grid2>
        </>
    );
};

export default LoginPage;
