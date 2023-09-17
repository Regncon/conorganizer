'use client';

import Link from 'next/link';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { Box, Card } from '../lib/mui';
import EventHeader from './eventHeader';

const EventList = () => {
    const { events, loading } = useAllEvents();

    return (
        <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20">
            {loading ? <h1>Loading...</h1> : null}
            {events?.map((conEvent) => (
                <Card key={conEvent.id} component={Link} href={`/event/${conEvent.id}`}>
                    <EventHeader conEvent={conEvent} />
                </Card>
            ))}
        </Box>
    );
};

export default EventList;
