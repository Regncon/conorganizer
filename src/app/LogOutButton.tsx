'use client';
import { signOutAndDeleteCookie } from '$lib/firebase/firebase';
import Button from '@mui/material/Button';
import { redirect, useRouter } from 'next/navigation';

const LogOutButton = () => {
    const router = useRouter();
    return (
        <Button
            onClick={async () => {
                await signOutAndDeleteCookie();
                redirect('/login');
            }}
            variant="contained"
        >
            logg ut
        </Button>
    );
};

export default LogOutButton;
