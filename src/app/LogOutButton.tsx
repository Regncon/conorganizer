'use client';
import { signOutAndDeleteCookie } from '$lib/firebase/firebase';
import { faSpinner } from '@fortawesome/free-solid-svg-icons/faSpinner';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Button from '@mui/material/Button';
import { redirect, useRouter } from 'next/navigation';
import { useTransition } from 'react';

const LogOutButton = () => {
    const router = useRouter();
    const [isPending, startTransition] = useTransition();
    return (
        <Button
            onClick={async () => {
                startTransition(async () => {
                    await signOutAndDeleteCookie();
                    router.replace('/login');
                });
            }}
            variant="contained"
            disabled={isPending}
            endIcon={isPending ? <FontAwesomeIcon icon={faSpinner} spin /> : undefined}
        >
            logg ut
        </Button>
    );
};

export default LogOutButton;
