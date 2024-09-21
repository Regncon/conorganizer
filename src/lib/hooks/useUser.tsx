import { firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged, onIdTokenChanged, type User } from 'firebase/auth';
import { useEffect, useRef, useState } from 'react';

const useUser = (checkForEmailVerification = false) => {
    const [isVerified, setIsVerified] = useState<boolean | undefined>(undefined);
    const userRef = useRef<User | null>(null);
    const fiveSeconds = 5_000;

    useEffect(() => {
        let userReloadInterval: NodeJS.Timeout | undefined;
        const unsubscribeUser = onAuthStateChanged(firebaseAuth, async (user) => {
            userRef.current = user;
            if (checkForEmailVerification) {
                userReloadInterval = setInterval(() => {
                    if (!isVerified) {
                        user?.reload().then(() => {
                            userRef.current = user;
                            setIsVerified(user?.emailVerified);
                        });
                    }
                }, fiveSeconds);
            }
        });

        return () => {
            unsubscribeUser();
            clearInterval(userReloadInterval);
        };
    }, [checkForEmailVerification, isVerified]);

    return { user: userRef.current, isUserVerified: isVerified };
};

export default useUser;
