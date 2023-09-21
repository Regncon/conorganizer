'use client';

import { useEffect, useState } from 'react';
import AccountCircle from '@mui/icons-material/AccountCircle';
import FilterAlt from '@mui/icons-material/FilterAlt';
import { Route } from 'next';
import Link from 'next/link';
import { Pool } from '@/lib/enums';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { useUserSettings } from '@/lib/hooks/useUserSettings';
import { Box, Card, Chip } from '../lib/mui';
import { useAuth } from './AuthProvider';
import EventHeader from './eventHeader';
import DayTab from './poolSelector';

const EventList = () => {
    const { events, loading } = useAllEvents();
    const [displayPool, setDisplayPool] = useState<Pool>(Pool.FridayEvening);
    console.log('displayPool', displayPool)
    const [showFilters, setShowFilters] = useState(false);
    const [showUnpublished, setShowUnpublished] = useState(false);

    const user = useAuth();
    const { conAuthorization } = useUserSettings(user?.uid);

    useEffect(() => {
        setShowUnpublished(conAuthorization?.admin && user ? true : false);
    }, [user, conAuthorization]);


    return (
        <>
            <DayTab handlePoolChange={(pool) => setDisplayPool(pool)} />
            <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20 mt-20">
                {loading ? <h1>Loading...</h1> : null}
                <Box sx={{ display: 'flex', gap: '.5em', flexGrow: '1', justifyContent: 'center', width: '100%' }}>
                    <Chip label="Alle" />
                    <Chip label="Mine p&aring;meldinger" variant="outlined" icon={<AccountCircle />} />
                    <Chip
                        icon={<FilterAlt />}
                        label="Andre filtre"
                        variant={showFilters ? 'filled' : 'outlined'}
                        onClick={() => setShowFilters(!showFilters)}
                    />
                </Box>
                <Box
                    display={showFilters ? 'flex' : 'none'}
                    sx={{ gap: '.5em', flexGrow: '1', justifyContent: 'center', width: '100%' }}
                >
                    <Chip label="Barnevennlig" variant="outlined" />
                    <Chip label="Rollespill" variant="outlined" />
                    <Chip label="Brettspill" variant="outlined" />
                    <Chip label="Annet" variant="outlined" />
                </Box>
                {events?.filter((conEvent) => conEvent.pool === displayPool)
                .filter((conEvent) => showUnpublished || conEvent.published)
                .map(
                    (
                        conEvent 
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
