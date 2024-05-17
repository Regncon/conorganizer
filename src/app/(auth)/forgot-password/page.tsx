'use client';
import Paper from '@mui/material/Paper';
import EmailTextField from '../shared/EmailTextField';
import Button from '@mui/material/Button';
import { useRouter, useSearchParams } from 'next/navigation';
import { forgotPassword } from '$lib/firebase/firebase';
import { useState, useTransition, type ComponentProps } from 'react';
import Container from '@mui/material/Container';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { Typography } from '@mui/material';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSpinner } from '@fortawesome/free-solid-svg-icons/faSpinner';

const ForgotPassword = () => {
    const [isPending, startTransition] = useTransition();
    const [message, setMessage] = useState<string>('');
    const router = useRouter();
    const searchParams = useSearchParams();
    const searchParamEmail = searchParams.get('email') ?? undefined;
    const handleSubmit: ComponentProps<'form'>['onSubmit'] = async (e) => {
        e.preventDefault();
        const { email } = Object.fromEntries(new FormData(e.target as HTMLFormElement)) as { email: string };
        setMessage(
            'Ein lenkje for å tilbakestille passordet er sendt viss du har ei registrert e-postadresse hos oss. Du vil bli omdirigert til innloggingssida.'
        );
        startTransition(async () => {
            await forgotPassword(email);
            await new Promise((resolve) => {
                setTimeout(async () => {
                    resolve(router.replace('/login'));
                }, 4000);
            });
        });
    };
    return (
        <Container component={Paper} fixed maxWidth="xl" sx={{ height: '70dvh' }}>
            <Grid2
                container
                component="form"
                sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
                onSubmit={handleSubmit}
            >
                <Typography>{message}</Typography>
                <EmailTextField defaultValue={searchParamEmail ?? undefined} />
                <Button
                    type="submit"
                    disabled={isPending}
                    endIcon={isPending ? <FontAwesomeIcon icon={faSpinner} spin /> : undefined}
                >
                    Gløymd passord?
                </Button>
            </Grid2>
        </Container>
    );
};

export default ForgotPassword;
