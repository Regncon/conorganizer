'use client';
import { Button, CircularProgress } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { singUpAndCreateCookie, type RegisterDetails } from '$lib/firebase/firebase';
import PasswordTextField from '../login/PasswordTextField';
import { useEffect, useTransition } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { emailRegExp, updateSearchParamsWithEmail } from '../shared/utils';
import EmailTextField from '../shared/EmailTextField';
import { validateRegisterFormAction } from './actions';
import { useFormState } from 'react-dom';

export const initialRegisterFormState = {
    emailError: '',
    passwordError: '',
    confirmError: '',
};
export type InitialRegisterFormState = typeof initialRegisterFormState;

const RegisterPage = () => {
    const searchParams = useSearchParams();
    const email = searchParams.get('email') ?? '';

    const [isPending, startTransition] = useTransition();
    const router = useRouter();
    useEffect(() => {
        router.prefetch('/dashboard');
    }, []);

    const handleFormSubmit: (formData: FormData) => Promise<void> = async (formData) => {
        const { password, confirm, email } = Object.fromEntries(formData) as RegisterDetails;

        if (password !== confirm) {
            console.log('no match', password !== confirm);
        }

        if (emailRegExp.test(email)) {
            startTransition(async () => {
                await singUpAndCreateCookie(formData);
                router.push('/dashboard');
            });
        }
    };

    const validateRegisterForm: (
        state: InitialRegisterFormState,
        formData: FormData
    ) => Promise<InitialRegisterFormState> = async (_, formData) => {
        const validatedState = await validateRegisterFormAction(formData);
        if (!validatedState.emailError && !validatedState.passwordError && !validatedState.confirmError) {
            await handleFormSubmit(formData);
        }
        return validatedState;
    };

    const [state, formAction] = useFormState(validateRegisterForm, initialRegisterFormState);

    return (
        <Grid2
            component="form"
            container
            sx={{ placeContent: 'center', flexDirection: 'column', minWidth: '20rem', gap: '1rem' }}
            onChange={(e) => {
                updateSearchParamsWithEmail(e, router, '/register');
            }}
            action={formAction}
        >
            <EmailTextField defaultValue={email} error={!!state.emailError} helperText={state.emailError} />
            <PasswordTextField
                autoComplete="new-password"
                error={!!state.passwordError}
                helperText={state.passwordError}
            />
            <PasswordTextField
                autoComplete="new-password"
                label="bekreft passord"
                name="confirm"
                error={!!state.confirmError}
                helperText={state.confirmError}
            />
            <Button
                fullWidth
                type="submit"
                disabled={isPending}
                endIcon={isPending ? <CircularProgress size="1.5rem" /> : undefined}
            >
                Lag ny brukar
            </Button>
        </Grid2>
    );
};
export default RegisterPage;
