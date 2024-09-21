'use client';
import GoogleSignInButton from '$app/(auth)/login/components/ui/GoogleButton';
import useUser from '$lib/hooks/useUser';
import { Button, CircularProgress } from '@mui/material';
import { sendEmailVerification, type User } from 'firebase/auth';
import { useState } from 'react';

type Props = {
    disabled: boolean | undefined;
};

const ConfirmEmailButtons = ({ disabled = false }: Props) => {
    const [isVerificationStarted, setIsVerificationStarted] = useState<boolean>(false);
    const { user, isUserVerified } = useUser(isVerificationStarted);
    const handleClick = async () => {
        setIsVerificationStarted(true);
        if (user && !disabled) {
            await sendEmailVerification(user);
        }
    };

    return (
        <>
            <Button
                fullWidth
                variant="contained"
                color="primary"
                onClick={handleClick}
                disabled={disabled || isUserVerified || isVerificationStarted}
            >
                {disabled || isUserVerified ?
                    'allerede verifisert'
                : isVerificationStarted ?
                    'Venter på verifisering av epost'
                :   'Bekreft e-post'}
                {!disabled && !isUserVerified && isVerificationStarted ?
                    <CircularProgress color="secondary" size="1.5rem" sx={{ marginInlineStart: '1rem' }} />
                :   null}
            </Button>
            {disabled || isUserVerified ? null : (
                <GoogleSignInButton redirectTo="/my-profile/my-tickets" disabled={isVerificationStarted} />
            )}
        </>
    );
};

export default ConfirmEmailButtons;
