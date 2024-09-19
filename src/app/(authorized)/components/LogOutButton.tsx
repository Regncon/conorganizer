'use client';
import { signOutAndDeleteCookie } from '$lib/firebase/firebase';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import { useRouter } from 'next/navigation';
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
            color="secondary"
            variant="contained"
            disabled={isPending}
            endIcon={isPending ? <CircularProgress size="1.5rem" /> : undefined}
        >
            logg ut
        </Button>
    );
};

export default LogOutButton;
