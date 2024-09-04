'use client';
import { db } from '$lib/firebase/firebase';
import { collection, onSnapshot } from 'firebase/firestore';
import { useEffect } from 'react';
import { updateDashboardEvents, updateEvents } from './serverAction';
type Props = {
    DashboardEvents: boolean;
};
const RealtimeEvents = ({ DashboardEvents = false }: Props) => {
    useEffect(() => {
        const eventsRef = collection(db, 'events');
        const unsubscribe = onSnapshot(eventsRef, (snapshot) => {
            updateDashboardEvents ? updateDashboardEvents() : updateEvents();
        });

        return () => {
            unsubscribe();
        };
    }, []);
    return null;
};

export default RealtimeEvents;
