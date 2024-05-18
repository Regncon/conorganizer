'use client';
import { Button, Container, Paper } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { singUpAndCreateCookie, type RegisterDetails } from '$lib/firebase/firebase';
import PasswordTextField from '../login/PasswordTextField';
import { useEffect, useRef, useTransition, type ComponentProps } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { emailRegExp, updateSearchParamsWithEmail } from '../shared/utils';
import EmailTextField from '../shared/EmailTextField';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSpinner } from '@fortawesome/free-solid-svg-icons/faSpinner';

const Register = () => {
    const searchParams = useSearchParams();
    const email = searchParams.get('email') ?? '';

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
        <Grid2
            component="form"
            container
            sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
            onChange={(e) => {
                updateSearchParamsWithEmail(e, router, '/register');
            }}
            onSubmit={handleSubmit}
        >
            <EmailTextField defaultValue={email} />
            <PasswordTextField autoComplete="new-password" />
            <PasswordTextField autoComplete="new-password" label="bekreft passord" name="confirm" />
            <Button
                type="submit"
                disabled={isPending}
                endIcon={isPending ? <FontAwesomeIcon icon={faSpinner} spin /> : undefined}
            >
                Log inn
            </Button>
        </Grid2>
    );
};

export default Register;
