'use client';
import { Button, Container, Link, Paper } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import PasswordTextField from './PasswordTextField';
import { forgotPassword, signInAndCreateCookie } from '$lib/firebase/firebase';
import { useEffect, useTransition, type ComponentProps } from 'react';
import EmailField from '../shared/ui/EmailField';
import { useRouter, useSearchParams } from 'next/navigation';
import type { Route } from 'next';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSpinner } from '@fortawesome/free-solid-svg-icons/faSpinner';

const Login = () => {
    const [isPending, startTransition] = useTransition();

    const router = useRouter();
    const searchParams = useSearchParams();
    const email = searchParams.get('email') ?? '';

    useEffect(() => {
        router.prefetch('/dashboard');
    }, []);

    const handleFormSubmit: ComponentProps<'form'>['onSubmit'] = async (e) => {
        startTransition(async () => {
            await signInAndCreateCookie(e);
            router.replace('/dashboard');
        });
    };

    const handleFormChange: ComponentProps<'form'>['onChange'] = async (e) => {
        const { value, name } = e.target as HTMLInputElement;
        if (name === 'email') {
            router.replace(`${'/login'}?email=${value}` as Route);
        }
    };

    const handleForgotPasswordClick: ComponentProps<'button'>['onClick'] = async (e) => {
        startTransition(() => {
            router.push(`/forgot-password?email=${email}`);
        });
    };

    const handleRegisterNewUser: ComponentProps<'button'>['onClick'] = async (e) => {
        startTransition(() => {
            router.push(`/register`);
        });
    };

    const disableAndLoadingSpinner = {
        disabled: isPending,
        endIcon: isPending ? <FontAwesomeIcon icon={faSpinner} spin /> : undefined,
    };

    return (
        <Container component={Paper} fixed maxWidth="xs" sx={{ height: '70dvh' }}>
            <Grid2
                component="form"
                container
                sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
                onSubmit={handleFormSubmit}
                onChange={handleFormChange}
            >
                <EmailField />
                <PasswordTextField />
                <Button type="submit" {...disableAndLoadingSpinner}>
                    Logg inn
                </Button>
                <Button onClick={handleForgotPasswordClick} {...disableAndLoadingSpinner}>
                    Gløymd passord?
                </Button>
                <Button
                    fullWidth
                    sx={{ marginLeft: 'auto', marginRight: 'auto' }}
                    onClick={handleRegisterNewUser}
                    {...disableAndLoadingSpinner}
                >
                    Registrer ny brukar
                </Button>
            </Grid2>
        </Container>
    );
};

export default Login;
