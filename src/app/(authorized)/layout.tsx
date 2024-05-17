import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import Box from '@mui/material/Box';
import { redirect } from 'next/navigation';
import type { ReactNode } from 'react';
import BackButton from './BackButton';
import LogOutButton from './LogOutButton';

type Props = {
    children: ReactNode;
};

const Layout = async ({ children }: Props) => {
    const { auth } = await getAuthorizedAuth();
    if (auth) {
        return (
            <>
                <Box sx={{ display: 'flex', placeContent: 'space-between', marginBlockStart: '1rem' }}>
                    <BackButton />
                    <LogOutButton />
                </Box>
                {children}
            </>
        );
    }
    redirect('/login');
};

export default Layout;
