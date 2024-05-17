'use client';
import { Button, Container, Paper } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { singUpAndCreateCookie, type RegisterDetails } from '$lib/firebase/firebase';
import PasswordTextField from '../login/PasswordTextField';
import { useEffect, useRef, useTransition, type ComponentProps } from 'react';
import { useRouter } from 'next/navigation';
import { emailRegExp } from '../shared/utils';
import EmailField from '../shared/ui/EmailField';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSpinner } from '@fortawesome/free-solid-svg-icons/faSpinner';

const Register = () => {
    const passwordRef = useRef<HTMLInputElement>(null);
    const confirmPasswordRef = useRef<HTMLInputElement>(null);

    const [isPending, startTransition] = useTransition();
    const router = useRouter();
    useEffect(() => {
        router.prefetch('/dashboard');
    }, []);

    const handleSubmit: ComponentProps<'form'>['onClick'] = async (e) => {
        e.preventDefault();
        const { password, confirm, email } = Object.fromEntries(
            new FormData(e.target as HTMLFormElement)
        ) as RegisterDetails;

        if (password !== confirm) {
            console.log('no match', password !== confirm);
        }

        if (emailRegExp.test(email)) {
            startTransition(async () => {
                await singUpAndCreateCookie(e);
                router.push('/dashboard');
            });
        }
    };

    return (
        <Container component={Paper} fixed maxWidth="xl" sx={{ height: '70dvh' }}>
            <Grid2
                component="form"
                container
                sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
                onSubmit={handleSubmit}
            >
                <EmailField />
                <PasswordTextField autoComplete="new-password" ref={passwordRef} />
                <PasswordTextField
                    autoComplete="new-password"
                    label="bekreft passord"
                    name="confirm"
                    ref={confirmPasswordRef}
                />
                <Button
                    type="submit"
                    disabled={isPending}
                    endIcon={isPending ? <FontAwesomeIcon icon={faSpinner} spin /> : undefined}
                >
                    Log inn
                </Button>
            </Grid2>
        </Container>
    );
};

export default Register;
