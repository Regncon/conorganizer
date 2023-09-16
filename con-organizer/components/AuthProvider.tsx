'use client';
import { createContext, useEffect, useState } from 'react';
import { onAuthStateChanged, User } from 'firebase/auth';
import { auth } from '../lib/firebase';
export const AuthContext = createContext<Partial<User | null>>({});
type Props = {
    children: JSX.Element;
};

export const AuthProvider = ({ children }: Props) => {
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        onAuthStateChanged(auth, (user) => {
            setCurrentUser(user);
            setLoading(false);
            console.log('user', user);
        });
    }, []);

    if (loading) {
        return (
            <div
                style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    height: '80vh',
                }}
            >
                <h1>Loading User...</h1>
            </div>
        );
    }

    return <AuthContext.Provider value={currentUser}>{children}</AuthContext.Provider>;
};
