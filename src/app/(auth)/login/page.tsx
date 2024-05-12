'use client';
import { Button, Container, Paper } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import PasswordTextField from './PasswordTextField';
import { signInAndCreateCookie, signOutAndDeleteCookie } from '$lib/firebase/firebase';

import { useRouter } from 'next/router';
import type { FormEvent } from 'react';
import EmailField from '../shared/ui/EmailField';

const Login = () => {
    const router = useRouter();

    const handleClick = async (e: FormEvent<HTMLFormElement>) => {
        await signInAndCreateCookie(e);
        router.push('/dashboard');
    };
    return (
        <Container component={Paper} fixed maxWidth="xl" sx={{ height: '70dvh' }}>
            <Grid2
                component="form"
                container
                sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
                onSubmit={handleClick}
            >
                <EmailField />
                <PasswordTextField />
                <Button type="submit">Log inn</Button>
            </Grid2>
            <Button onClick={signOutAndDeleteCookie}>logg ut</Button>
        </Container>
    );
};

export default Login;
