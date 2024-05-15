'use client';
import { Button, Container, Link, Paper } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import PasswordTextField from './PasswordTextField';
import { forgotPassword, signInAndCreateCookie } from '$lib/firebase/firebase';
import type { ComponentProps } from 'react';
import EmailField from '../shared/ui/EmailField';
import { useRouter, useSearchParams } from 'next/navigation';
import type { Route } from 'next';

const Login = () => {
    const router = useRouter();
    const searchParams = useSearchParams();
    const email = searchParams.get('email');

    const handleFormSubmit: ComponentProps<'form'>['onSubmit'] = async (e) => {
        await signInAndCreateCookie(e);
        router.push('/dashboard');
    };
    const handleFormChange: ComponentProps<'form'>['onChange'] = async (e) => {
        const { value, name } = e.target as HTMLInputElement;
        if (name === 'email') {
            router.push(`${'/login'}?email=${value}` as Route);
        }
    };
    const handleForgotPasswordClick: ComponentProps<'button'>['onClick'] = async (e) => {
        router.push(`/forgot-password?email=${email}`);
    };
    return (
        <Container component={Paper} fixed maxWidth="xl" sx={{ height: '70dvh' }}>
            <Grid2
                component="form"
                container
                sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
                onSubmit={handleFormSubmit}
                onChange={handleFormChange}
            >
                <EmailField />
                <PasswordTextField />
                <Button type="submit">Logg inn</Button>
                <Button onClick={handleForgotPasswordClick}>Glemt passord?</Button>
                <Link sx={{ marginLeft: 'auto', marginRight: 'auto' }} href="/register">
                    Registrer ny bruker
                </Link>
            </Grid2>
        </Container>
    );
};

export default Login;
