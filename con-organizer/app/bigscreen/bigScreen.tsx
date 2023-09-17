'use client';

import { useEffect, useState } from 'react';
import { Box, Card, CardContent, CardHeader, Divider } from '@mui/material';
import { collection, onSnapshot } from 'firebase/firestore';
import EventHeader from '@/components/eventHeader';
import { pool } from '@/lib/enums';
import { ConEvent } from '@/lib/types';
import db from '../../lib/firebase';

const BigScreen = () => {
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
        <Box className="flex flex-row flex-wrap justify-center gap-4">
            <Box className="flex flex-col gap-4 bg-slate-800 p-4"
            sx={{ maxWidth: '440px' }}>
                <h1>Fredag Kveld</h1>
                <p>18:00 - 23:00</p>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Innsjekk Fredag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 16:30 - 18:00 </p>
                    </CardContent>
                </Card>
                {conEvents
                    .filter((conEvent) => conEvent.pool === pool.FirdayEvening)
                    .map((conEvent) => (
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4 bg-slate-900 p-4"
            sx={{ maxWidth: '440px' }}>
                <h1>Lørdag Morgen </h1>
                <p>10:00 - 15:00</p>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Innsjekk Lørdag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 09:00 - 10:00 </p>
                    </CardContent>
                </Card>
                {conEvents
                    .filter((conEvent) => conEvent.pool === pool.SaturdayMorning)
                    .map((conEvent) => (
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4 bg-slate-800 p-4"
            sx={{ maxWidth: '440px' }}>
                <h1>Lørdag Kveld </h1>
                <p>18:00 - 23:00</p>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Middag Lørdag" />

                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 16:00 - 18:00 </p>
                        <p>Påmeling til middag er frem til 4. oktober for eksempel?</p>
                    </CardContent>
                </Card>
                {conEvents
                    .filter((conEvent) => conEvent.pool === pool.SaturdayEvening)
                    .map((conEvent) => (
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4 bg-slate-900 p-4"
            sx={{ maxWidth: '440px' }}>
                <h1>Søndag Morgen </h1>
                <p>10:00 - 15:00</p>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Innsjekk Søndag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 09:00 - 10:00 </p>
                    </CardContent>
                </Card>
                {conEvents
                    .filter((conEvent) => conEvent.pool === pool.SundayMorning)
                    .map((conEvent) => (
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    ))}

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Prisutdeling" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 16:00 - 17:00 </p>
                    </CardContent>
                </Card>
            </Box>
        </Box>
    );
};
export default BigScreen;
