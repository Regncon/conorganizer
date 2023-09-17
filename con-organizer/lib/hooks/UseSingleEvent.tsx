import { useEffect, useState } from 'react';
import { ConEvent } from '../../models/types';
import { singleEvent$ } from '../observables/SingleEvent';

export const useSingleEvents = (id: string) => {
    const [event, setEvent] = useState<ConEvent>();
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const eventObservable = singleEvent$(id).subscribe((snapshot) => {
            if (snapshot.data()) {
                setEvent({ ...(snapshot.data() as ConEvent), id: snapshot.id });
                setLoading(false);
            }
            setLoading(false);
        });

        return () => {
            eventObservable.unsubscribe();
        };
    }, []);

    return { event, loading };
};
