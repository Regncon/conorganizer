'use client';
import { firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged } from 'firebase/auth';
import { useEffect } from 'react';
import { setSessionCookie } from '../(auth)/login/lib/actions';

type Props = {};

const SetCookie = ({}: Props) => {
    useEffect(() => {
        const unsubscribeUser = onAuthStateChanged(firebaseAuth, async (user) => {
            console.log(user, 'user');

            if (!user) {
                return;
            }
            const idToken = await user.getIdToken();
            console.log(idToken, 'idToken');
            //ToDo: only set cookie if it is not already set or the cookie is expired
            //console.log('setting cookie');
            await setSessionCookie(idToken);
        });

        return () => {
            unsubscribeUser();
        };
    }, []);
    return null;
};

export default SetCookie;
