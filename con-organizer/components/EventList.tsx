'use client';

import { useEffect, useState } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import AccountCircle from '@mui/icons-material/AccountCircle';
import FilterAlt from '@mui/icons-material/FilterAlt';
import { Pool } from '@/lib/enums';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { useUserSettings } from '@/lib/hooks/UseUserSettings';
import { Box, Chip, Typography } from '../lib/mui';
import EventCardBoundary from './ErrorBoundaries/EventCardBoundary';
import { useAuth } from './AuthProvider';
import EventCard from './EventCard';
import PoolSelector from './PoolSelector';

const EventList = () => {
    const { events, loading } = useAllEvents();
    const [displayPool, setDisplayPool] = useState<Pool>(Pool.FridayEvening);
    const [showFilters, setShowFilters] = useState(false);
    const [showUnpublished, setShowUnpublished] = useState(false);

    const user = useAuth();
    const { userSettings } = useUserSettings(user?.uid);

    useEffect(() => {
        setShowUnpublished(userSettings?.admin && user ? true : false);
    }, [user, userSettings]);

    return (
        <>
            <PoolSelector handlePoolChange={(pool) => setDisplayPool(pool)} />
            <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20 mt-20">
                {loading ? <Typography variant="body1">Loading...</Typography> : null}
                {/* <Box sx={{ display: 'flex', gap: '.5em', flexGrow: '1', justifyContent: 'center', width: '100%' }}>
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
                </Box> */}
                {events
                    ?.filter((conEvent) => conEvent.pool === displayPool)
                    .filter((conEvent) => showUnpublished || conEvent.published)
                    .map((conEvent) => (
                        <ErrorBoundary FallbackComponent={EventCardBoundary} key={conEvent.id}>
                            <EventCard conEvent={conEvent} />
                        </ErrorBoundary>
                    ))}
            </Box>
        </>
    );
};

export default EventList;
