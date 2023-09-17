'use client';

import Link from 'next/link';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { Box, Card, CardContent, CardHeader } from '../lib/mui';
import AddEvent from './addEvent';
import EventHeader from './eventHeader';

const EventList = () => {
    const { events, loading } = useAllEvents();

    return (
        <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20 mt-20">
            {loading ? <h1>Loading...</h1> : null}
            <AddEvent />
            <Card sx={{ width: '100%' }}>
                <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Registrering Fredag" />
                <CardContent sx={{ paddingTop: '0' }}>
                    <p>Kl 16:00 - 17:00 </p>
                </CardContent>
            </Card>
            {events?.map(
                (
                    conEvent //filter((conEvent) => conEvent.published === true)
                ) => (
                    <Card
                        key={conEvent.id}
                        component={Link}
                        href={`/event/${conEvent.id}`}
                        sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                    >
                        <EventHeader conEvent={conEvent} />
                    </Card>
                )
            )}
        </Box>
    );
};

export default EventList;
