'use client';
import useUser from '$lib/hooks/useUser';
import { Button, CircularProgress } from '@mui/material';
import { sendEmailVerification, type User } from 'firebase/auth';

type Props = {};

const ConfirmEmailButton = ({}: Props) => {
    const user = useUser();
    const handleClick = async () => {
        if (user && !user.emailVerified) {
            await sendEmailVerification(user);
        }
    };
    return (
        <Button fullWidth variant="contained" color="primary" onClick={handleClick} disabled={!user}>
            Bekreft e-post {!user && <CircularProgress size="1.5rem" sx={{ marginInlineStart: '1rem' }} />}
        </Button>
    );
};

export default ConfirmEmailButton;
