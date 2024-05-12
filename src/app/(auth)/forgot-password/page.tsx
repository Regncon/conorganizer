'use client';
import Paper from '@mui/material/Paper';
import EmailField from '../shared/ui/EmailField';
import Button from '@mui/material/Button';
import { useRouter, useSearchParams } from 'next/navigation';
import { forgotPassword } from '$lib/firebase/firebase';
import type { ComponentProps } from 'react';
import Container from '@mui/material/Container';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';

const ForgotPassword = () => {
    const router = useRouter();
    const searchParams = useSearchParams();
    const searchParamEmail = searchParams.get('email') ?? undefined;
    const handleSubmit: ComponentProps<'form'>['onSubmit'] = async (e) => {
        e.preventDefault();
        const { email } = Object.fromEntries(new FormData(e.target as HTMLFormElement)) as { email: string };
        await forgotPassword(email);
        router.replace('/login');
    };
    return (
        <Container component={Paper} fixed maxWidth="xl" sx={{ height: '70dvh' }}>
            <Grid2
                container
                component="form"
                sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
                onSubmit={handleSubmit}
            >
                <EmailField defaultValue={searchParamEmail ?? undefined} />
                <Button type="submit">Glemt passord?</Button>
            </Grid2>
        </Container>
    );
};

export default ForgotPassword;
