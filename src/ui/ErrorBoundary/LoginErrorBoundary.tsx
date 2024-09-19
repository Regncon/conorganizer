'use client'; // Error components must be Client Components

import { Box, Button, Typography } from '@mui/material';
import Link from 'next/link';
import type { PropsWithChildren } from 'react';

const FirebaseLoginError = {
    InvalidEmail: 'auth/invalid-email',
    InvalidCredential: 'auth/invalid-credential',
    EmailAlreadyInUse: 'auth/email-already-in-use',
    TooManyRequests: 'auth/too-many-requests',
    WrongPassword: 'auth/wrong-password',
} as const;

const LoginErrorBoundary: ErrorBoundaryProps = ({ error, reset }) => {
    const ResetButton = (
        <Button onClick={reset} variant="contained">
            Prøv igjen
        </Button>
    );
    const InlineWrapper = ({ children, marginBottom }: { marginBottom?: boolean } & PropsWithChildren) => (
        <Box display="flex" flexDirection="row" marginBlockEnd={marginBottom ? '1rem' : 'unset'}>
            {children}
        </Box>
    );

    if (error.message.includes(FirebaseLoginError.InvalidEmail)) {
        return (
            <>
                <Typography variant="h2">
                    Det ser ut som du skreiv noko feil i e-postadressa di, vennlegast prøv igjen.
                </Typography>
                {ResetButton}
            </>
        );
    }
    if (error.message.includes(FirebaseLoginError.InvalidCredential)) {
        return (
            <>
                <Typography variant="h2">Passordet eller e-postadressa er ikkje korrekt.</Typography>
                {ResetButton}
            </>
        );
    }
    if (error.message.includes(FirebaseLoginError.EmailAlreadyInUse)) {
        return (
            <>
                <Typography variant="h2">Det ser ut til at e-postadressa allereie er registrert hos oss.</Typography>
                <InlineWrapper marginBottom>
                    <Typography variant="h3" margin="0" marginInlineEnd="1rem">
                        Gå til
                    </Typography>
                    <Button component={Link} href="/login" variant="contained">
                        Log inn
                    </Button>
                </InlineWrapper>

                <InlineWrapper>
                    <Typography variant="h3" margin="0" marginInlineEnd="1rem">
                        Eller
                    </Typography>
                    {ResetButton}
                </InlineWrapper>
            </>
        );
    }
    if (error.message.includes(FirebaseLoginError.TooManyRequests)) {
        return (
            <>
                <Typography variant="h2">
                    Vi har mellombels sperra kontoen din på grunn av for mange påloggingsforsøk. Dette er for å hindre
                    at hackarar og botar kan gjette seg til passordet ditt. Vent litt før du prøver igjen.
                </Typography>
                {ResetButton}
            </>
        );
    }
    if (error.message.includes(FirebaseLoginError.WrongPassword)) {
        return (
            <>
                <Typography variant="h2">
                    Det ser ut som du har skrive noko feil i passordet ditt, vennlegast prøv igjen.
                </Typography>
                ;{ResetButton}
            </>
        );
    }

    return (
        <>
            <Typography variant="h3">
                {`Kunne ikkje logge deg inn fordi det oppstod ein feil. Tekniske detaljar: ${error.message}`}
            </Typography>
            {ResetButton}
        </>
    );
};
export default LoginErrorBoundary;
