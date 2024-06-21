'use client';
import { Button, CircularProgress, Typography } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import PasswordTextField from './PasswordTextField';
import { signInAndCreateCookie } from '$lib/firebase/firebase';
import { useEffect, useState, useTransition, type ComponentProps } from 'react';
import EmailTextField from '../shared/EmailTextField';
import { useRouter, useSearchParams } from 'next/navigation';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSpinner } from '@fortawesome/free-solid-svg-icons/faSpinner';
import { useFormState, useFormStatus } from 'react-dom';
import { validateForm as validateFormAction } from './action';
import LoginButton from './LoginButton';
import { updateSearchParamsWithEmail } from '../shared/utils';

const initialFormState = {
    emailError: '',
    passwordError: '',
};
export type InitialFormState = typeof initialFormState;

export const disableAndLoadingSpinner = (shouldSpin: boolean, isPending: boolean) => ({
    disabled: isPending,
    endIcon: isPending && shouldSpin ? <CircularProgress size="1.5rem" /> : undefined,
});

const LoginPage = () => {
    const [isPending, startTransition] = useTransition();
    const [spinners, setSpinners] = useState<{ login: boolean; forgot: boolean; register: boolean }>({
        forgot: false,
        register: false,
    });

    const router = useRouter();
    const searchParams = useSearchParams();
    const email = searchParams.get('email') ?? '';
    const expiredSession = searchParams.get('expired') === 'true' ? true : false;

    useEffect(() => {
        router.prefetch('/dashboard');
    }, []);

    const handleFormSubmit = async (formData: FormData) => {
        startTransition(async () => {
            await signInAndCreateCookie(formData);
            router.replace('/dashboard');
        });
    };

    const handleForgotPasswordClick: ComponentProps<'button'>['onClick'] = async (e) => {
        setSpinners({ ...spinners, forgot: true });
        startTransition(() => {
            router.push(`/forgot-password?email=${email}`);
        });
    };

    const handleRegisterNewUser: ComponentProps<'button'>['onClick'] = async (e) => {
        setSpinners({ ...spinners, register: true });
        startTransition(() => {
            router.push(`/register`);
        });
    };

    const validateForm: (
        state: typeof initialFormState,
        formData: FormData
    ) => Promise<typeof initialFormState> = async (_, formData) => {
        const validatedState = await validateFormAction(formData);
        if (!validatedState.emailError && !validatedState.passwordError) {
            await handleFormSubmit(formData);
        }
        return validatedState;
    };

    const [state, formAction] = useFormState(validateForm, initialFormState);

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
                }}
                onChange={(e) => {
                    updateSearchParamsWithEmail(e, router, '/login');
                }}
                action={formAction}
            >
                <EmailTextField defaultValue={email} error={state.emailError} helperText={state.emailError} />
                <PasswordTextField error={state.passwordError} helperText={state.passwordError} />

                <LoginButton />

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
