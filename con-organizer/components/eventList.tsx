'use client';

import { useEffect, useState } from 'react';
import { collection, onSnapshot } from 'firebase/firestore';
import { ConEvent } from '@/lib/types';
import db from '../lib/firebase';
import { Box, Card, CardContent, CardHeader } from '../lib/mui';
import AddEvent from './addEvent';
import EventHeader from './eventHeader';

const EventList = () => {
    const collectionRef = collection(db, 'events');
    const [conEvents, setconEvents] = useState([] as ConEvent[]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);
        const unsub = onSnapshot(collectionRef, (querySnapshot) => {
            const items = [] as ConEvent[];
            querySnapshot.forEach((doc) => {
                items.push(doc.data() as ConEvent);
                items[items.length - 1].id = doc.id;
            });
            setconEvents(items);
            setLoading(false);
        });
        return () => {
            unsub();
        };
    }, []);

    return (
        <>
            <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20">
                {loading ? <h1>Loading...</h1> : null}
            <AddEvent collectionRef={collectionRef} />
            <Card sx={{ width: '100%' }}>
                <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Registrering Fredag" />
                <CardContent
                    sx={{ paddingTop: '0' }}
                >
                    <p>Kl 16:00 - 17:00 </p>
                </CardContent>
            </Card>
                {conEvents.map(
                    (
                        conEvent //filter((conEvent) => conEvent.published === true)
                    ) => (
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{cursor: 'pointer', opacity: conEvent?.published === false ? '50%'  : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    )
                )}
            </Box>
        </>
    );
};

export default EventList;
