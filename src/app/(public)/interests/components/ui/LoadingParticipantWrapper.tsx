'use client';
import { Box, CircularProgress, Typography } from '@mui/material';
import { useEffect, useState, type PropsWithChildren } from 'react';

type Props = {};

const LoadingParticipantWrapper = ({ children }: PropsWithChildren<Props>) => {
    const [isLoading, setIsLoading] = useState<boolean>(false);
    useEffect(() => {
        const handleLoading = (e: Event & { detail?: { loading?: boolean } }) => {
            setIsLoading(e.detail?.loading ?? false);
        };

        document.addEventListener('my-participants-changed', handleLoading);

        return () => {
            document.removeEventListener('my-participants-changed', handleLoading);
        };
    }, []);
    return isLoading ?
            <Box sx={{ display: 'grid', placeContent: 'center', alignItems: 'center', marginBlockStart: '2rem' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    <Typography>Venleg vent medan interesser vert lasta inn...</Typography> <CircularProgress />
                </Box>
            </Box>
        :   children;
};

export default LoadingParticipantWrapper;
