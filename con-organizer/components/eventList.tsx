'use client';

import { useState } from 'react';
import AccountCircle from '@mui/icons-material/AccountCircle';
import FilterAlt from '@mui/icons-material/FilterAlt';
import { Route } from 'next';
import Link from 'next/link';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { Box, Card, Chip } from '../lib/mui';
import AddEvent from './addEvent';
import EventHeader from './eventHeader';

const EventList = () => {
    const { events, loading } = useAllEvents();
    const [showFilters, setShowFilters] = useState(false);

    return (
        <>
            <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20 mt-20">
                {loading ? <h1>Loading...</h1> : null}
                <AddEvent />
                <Box display="flex" gap=".5em">
                    <Chip label="Alle" />
                    <Chip label="Mine p&aring;meldinger" variant="outlined" icon={<AccountCircle />} />
                    <Chip
                        icon={<FilterAlt />}
                        label="Andre filtre"
                        variant={showFilters ? 'filled' : 'outlined'}
                        onClick={() => setShowFilters(!showFilters)}
                    />
                </Box>
                <Box display={showFilters ? 'flex' : 'none'} gap=".5em">
                    <Chip label="Barnevennlig" variant="outlined" />
                    <Chip label="Rollespill" variant="outlined" />
                    <Chip label="Brettspill" variant="outlined" />
                    <Chip label="Annet" variant="outlined" />
                </Box>
                {events?.map(
                    (
                        conEvent //filter((conEvent) => conEvent.published === true)
                    ) => (
                        <Card
                            key={conEvent.id}
                            component={Link}
                            href={`/event/${conEvent.id}` as Route}
                            sx={{
                                maxWidth: '500px',
                                cursor: 'pointer',
                                opacity: conEvent?.published === false ? '50%' : '',
                            }}
                        >
                            <EventHeader conEvent={conEvent} listView={true} />
                        </Card>
                    )
                )}
            </Box>
        </>
    );
};

export default EventList;
