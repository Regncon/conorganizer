'use client';
import { updateEventById } from '$app/(public)/components/lib/serverAction';
import { db } from '$lib/firebase/firebase';
import type { Unsubscribe } from 'firebase/auth';
import { onSnapshot, doc } from 'firebase/firestore';
import { useEffect } from 'react';

type Props = {
    id: string | undefined;
};

const RealtimePoolEvent = ({ id }: Props) => {
    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (id !== undefined) {
            unsubscribeSnapshot = onSnapshot(doc(db, 'pool-events', id), (snapshot) => {
                console.log('here1', id, snapshot.data());

                updateEventById(id);
            });
        }
        return () => {
            unsubscribeSnapshot?.();
        };
    }, [id]);
    return null;
};

export default RealtimePoolEvent;
