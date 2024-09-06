import { getMyUserInfo } from '$app/(authorized)/my-events/actions';
import { redirect } from 'next/navigation';
import { getAuthorizedAuth } from './firebase/firebaseAdmin';

export const redirectToAdminDashboardWhenAdministrator = async () => {
    const { app, user, auth, db } = await getAuthorizedAuth();
    if (app !== null && user !== null && auth !== null && db !== null) {
        const userInfo = await getMyUserInfo(db, user);
        if (userInfo?.admin) {
            redirect('/admin/dashboard');
        }
    }
};
