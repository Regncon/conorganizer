'use client';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { collection, onSnapshot } from 'firebase/firestore';
import { useEffect } from 'react';
import { updateMyEvents } from '../lib/actions';
type Props = {
    userId: string;
};

const RealtimeMyEvents = ({ userId }: Props) => {
    useEffect(() => {
        const eventsRef = collection(db, 'users', userId, 'my-events');
        const unsubscribe = onSnapshot(eventsRef, () => {
            updateMyEvents();
        });

        return () => {
            unsubscribe();
        };
    }, []);
    return null;
};

export default RealtimeMyEvents;
