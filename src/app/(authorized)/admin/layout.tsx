import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Typography } from '@mui/material';
import Link from 'next/link';
import type { PropsWithChildren } from 'react';
import { getMyUserInfo } from '../my-events/lib/actions';

type Props = {} & PropsWithChildren;

const Layout = async ({ children }: Props) => {
    const { app, user, auth, db } = await getAuthorizedAuth();
    if (app !== null && user !== null && auth !== null && db !== null) {
        const userInfo = await getMyUserInfo(db, user);
        if (!userInfo?.admin) {
            return (
                <Typography variant="h5" component={Link} href="/dashboard" sx={{ marginTop: '1rem' }}>
                    Du er ikkje administrator, gå til oversikta
                </Typography>
            );
        }
    }
    return children;
};

export default Layout;
