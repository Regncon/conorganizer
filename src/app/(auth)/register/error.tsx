'use client'; // Error components must be Client Components

import { Button, Typography } from '@mui/material';
import Link from 'next/link';
import { useEffect } from 'react';

export default function Error({ error, reset }: { error: Error & { digest?: string }; reset: () => void }) {
    useEffect(() => {
        // Log the error to an error reporting service
        console.error(error);
    }, [error]);

    return (
        <>
            <h2>{error.message}</h2>
            <Button onClick={reset} variant="outlined">
                Prøv igjenn
            </Button>
        </>
    );
}
