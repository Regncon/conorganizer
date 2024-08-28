'use client';
import { firebaseAuth } from "$lib/firebase/firebase";
import { onAuthStateChanged } from "firebase/auth";
import { useEffect } from "react";
import { setSessionCookie } from "./(auth)/login/actions";

type Props = {}; 

const SetCookie = ({}:Props) => { 
    useEffect(() => {
        const unsubscribeUser = onAuthStateChanged(firebaseAuth, async (user) => {
            if (!user) {
                return;
            }
            const idToken = await user.getIdToken();

            await setSessionCookie(idToken);
        });

        return () => {
            unsubscribeUser();
        };
    }, []);
    return (null);
};

 export default SetCookie;