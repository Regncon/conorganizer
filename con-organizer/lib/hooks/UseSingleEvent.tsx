import { useEffect, useState } from 'react';
import { collection, query, where } from 'firebase/firestore';
import { collectionData } from 'rxfire/firestore';
import { tap } from 'rxjs';
import db from '@/lib/firebase';
import { ConEvent } from '../types';

export const useSingleEvents = (id: string) => {
    const [event, setEvent] = useState<ConEvent>();
    const [loading, setLoading] = useState<boolean>(true);
    const eventRef = query(collection(db, 'events'), where('id', '==', id));

    useEffect(() => {
        const conEvents = collectionData<ConEvent>(eventRef)
            .pipe(tap((event) => console.log('This is just an observable!')))
            .subscribe((event) => {
                setEvent(event[0]);
                setLoading(false);
            });

        return () => {
            conEvents.unsubscribe();
        };
    }, []);

    return { event, loading };
};
