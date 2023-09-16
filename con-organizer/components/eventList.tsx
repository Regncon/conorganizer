"use client";

import { Box, Card } from '../lib/mui';
import React, { useEffect, useState } from 'react';
import { onSnapshot, collection } from 'firebase/firestore';
import db from '../lib/firebase';
import { ConEvent } from '@/lib/types';
import EventHeader from './eventHeader';

// import parse from 'html-react-parser';
interface Props {}

const EventList = () => {
    const colletionRef = collection(db, 'schools');
    const [conEvents, setconEvents] = useState([] as ConEvent[]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);
        const unsub = onSnapshot(colletionRef, (querySnapshot) => {
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
        <Box className='flex flex-row flex-wrap justify-center gap-4 mb-20'>
            {loading ? <h1>Loading...</h1> : null}
            {conEvents.map((conEvent) => (

                <Card key={conEvent.id}
                    onClick={() => {
                        window.location.assign(`/event/${conEvent.id}`);
                    }}
                >
                    <EventHeader conEvent={conEvent} />
                </Card>

            ))}
        </Box>
    );
};

export default EventList;
