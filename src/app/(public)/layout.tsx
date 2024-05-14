import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Card, CardContent, Typography } from '@mui/material';
import Link from 'next/link';
import type { ReactNode } from 'react';

type Props = {
    children: ReactNode;
};

const Layout = async ({ children }: Props) => {
    const { auth } = await getAuthorizedAuth();
    return (
        <>
            {auth?.currentUser?.uid ? null : (
                <Card sx={{ marginTop: '1rem' }}>
                    <CardContent>
                        <Typography variant="h1">
                            For og lage arrangementer må du ha en bruker trykk på <Link href="/login">logginn</Link>{' '}
                            Eller
                            <Link href="/register"> registrer </Link>
                        </Typography>
                    </CardContent>
                </Card>
            )}
            {children}
        </>
    );
};

export default Layout;
