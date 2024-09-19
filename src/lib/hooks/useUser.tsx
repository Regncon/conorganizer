import { firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged, type User } from 'firebase/auth';
import { useEffect, useRef, useState } from 'react';

const useUser = () => {
    const userRef = useRef<User | null>(null);
    useEffect(() => {
        const unsubscribeUser = onAuthStateChanged(firebaseAuth, async (user) => {
            userRef.current = user;
            return;
        });

        return () => {
            unsubscribeUser();
        };
    }, []);

    return userRef.current;
};

export default useUser;
