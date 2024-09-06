'use client';
import { db } from '$lib/firebase/firebase';
import { collection, onSnapshot } from 'firebase/firestore';
import { useEffect } from 'react';
import { updateDashboardEvents, updateEvents } from './serverAction';
import { updateMyEvents } from '$app/(authorized)/my-events/lib/actions';

type Props = {
    where: 'DASHBOARD_EVENTS' | 'EVENTS';
};

const RealtimeEvents = ({ where = 'EVENTS' }: Props) => {
    useEffect(() => {
        const eventsRef = collection(db, 'events');
        const unsubscribe = onSnapshot(eventsRef, (snapshot) => {
            switch (where) {
                case 'DASHBOARD_EVENTS':
                    updateDashboardEvents();
                    break;
                case 'EVENTS':
                    updateEvents();
                    break;
            }
        });

        return () => {
            unsubscribe();
        };
    }, []);
    return null;
};

export default RealtimeEvents;
