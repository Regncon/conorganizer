'use client';
import { Button, Container, Link, Paper, Typography } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import PasswordTextField from './PasswordTextField';
import { forgotPassword, signInAndCreateCookie } from '$lib/firebase/firebase';
import { useEffect, useState, useTransition, type ComponentProps } from 'react';
import EmailTextField from '../shared/EmailTextField';
import { useRouter, useSearchParams } from 'next/navigation';
import type { Route } from 'next';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSpinner } from '@fortawesome/free-solid-svg-icons/faSpinner';
import { updateSearchParamsWithEmail } from '../shared/utils';

const LoginPage = () => {
    const [isPending, startTransition] = useTransition();
    const [spinners, setSpinners] = useState<{ login: boolean; forgot: boolean; register: boolean }>({
        forgot: false,
        login: false,
        register: false,
    });

    const router = useRouter();
    const searchParams = useSearchParams();
    const email = searchParams.get('email') ?? '';
    const expiredSession = searchParams.get('expired') === 'true' ? true : false;

    useEffect(() => {
        router.prefetch('/dashboard');
    }, []);

    const handleFormSubmit: ComponentProps<'form'>['onSubmit'] = async (e) => {
        setSpinners({ ...spinners, login: true });
        startTransition(async () => {
            await signInAndCreateCookie(e);
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

    const disableAndLoadingSpinner = (shouldSpin: boolean) => ({
        disabled: isPending,
        endIcon: isPending && shouldSpin ? <FontAwesomeIcon icon={faSpinner} spin /> : undefined,
    });

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
                onSubmit={handleFormSubmit}
                onChange={(e) => {
                    updateSearchParamsWithEmail(e, router, '/login');
                }}
            >
                <EmailTextField defaultValue={email} />
                <PasswordTextField />
                <Button type="submit" {...disableAndLoadingSpinner(spinners.login)}>
                    Logg inn
                </Button>
                <Button onClick={handleForgotPasswordClick} {...disableAndLoadingSpinner(spinners.forgot)}>
                    Gløymd passord?
                </Button>
                <Button
                    fullWidth
                    sx={{ marginLeft: 'auto', marginRight: 'auto' }}
                    onClick={handleRegisterNewUser}
                    {...disableAndLoadingSpinner(spinners.register)}
                >
                    Registrer ny brukar
                </Button>
            </Grid2>
        </>
    );
};

export default LoginPage;
