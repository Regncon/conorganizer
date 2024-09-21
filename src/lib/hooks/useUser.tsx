import { firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged, onIdTokenChanged, type User } from 'firebase/auth';
import { useEffect, useRef, useState } from 'react';

const useUser = (checkForEmailVerification = false) => {
    const [isVerified, setIsVerified] = useState<boolean | undefined>(undefined);
    const userRef = useRef<User | null>(null);
    const fiveSeconds = 5_000;

    useEffect(() => {
        const unsubscribeUser = onAuthStateChanged(firebaseAuth, async (user) => {
            userRef.current = user;
            if (checkForEmailVerification) {
                const userReloadInterval = setInterval(() => {
                    if (!isVerified) {
                        user?.reload().then(() => {
                            userRef.current = user;
                            setIsVerified(user?.emailVerified);
                        });
                        console.log(!isVerified, '!isVerified');
                    }
                }, fiveSeconds);
                if (isVerified) {
                    clearInterval(userReloadInterval);
                }
            }
        });

        return () => {
            unsubscribeUser();
        };
    }, [checkForEmailVerification, isVerified]);

    return { user: userRef.current, isUserVerified: isVerified };
};

export default useUser;
