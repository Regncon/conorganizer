'use client';

import { signOutAndDeleteCookie } from '$lib/firebase/firebase';
import { Box, CircularProgress, Paper, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect, startTransition } from 'react';

type Props = {};

const LogOut = ({ }: Props) => {
    const router = useRouter();
    useEffect(() => {
        startTransition(async () => {
            await signOutAndDeleteCookie();
            router.replace('/');
        });
    }, []);
    return (
        <Box sx={{ display: 'grid', placeContent: 'center', placeItems: 'center', height: '65%' }}>
            <Typography>Logger ut</Typography>
            <CircularProgress />
        </Box>
    );
};

export default LogOut;
