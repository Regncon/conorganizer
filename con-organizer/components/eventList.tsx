'use client';

import Link from 'next/link';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { Box, Card } from '../lib/mui';
import AddEvent from './addEvent';
import EventHeader from './eventHeader';

const EventList = () => {
    const { events, loading } = useAllEvents();

    return (
        <>
            <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20 mt-20">
                {loading ? <h1>Loading...</h1> : null}
                <AddEvent />
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
        </>
    );
};

export default EventList;
