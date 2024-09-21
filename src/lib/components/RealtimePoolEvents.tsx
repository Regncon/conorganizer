'use client';
import { updateEvents } from '$app/(public)/components/lib/serverAction';
import { db } from '$lib/firebase/firebase';
import { onSnapshot, collection } from 'firebase/firestore';
import { useEffect } from 'react';

type Props = {};

const RealtimePoolEvents = ({}: Props) => {
    useEffect(() => {
        const unsubscribeSnapshot = onSnapshot(collection(db, 'pool-events'), (_) => {
            updateEvents();
        });

        return () => {
            unsubscribeSnapshot?.();
        };
    }, []);
    return null;
};

export default RealtimePoolEvents;
