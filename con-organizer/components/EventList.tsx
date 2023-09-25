'use client';

import { useEffect, useState } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import AccountCircle from '@mui/icons-material/AccountCircle';
import FilterAlt from '@mui/icons-material/FilterAlt';
import { useSearchParams } from 'next/navigation';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { useUserSettings } from '@/lib/hooks/UseUserSettings';
import { Pool } from '@/models/enums';
import { Box, Chip, Typography } from '../lib/mui';
import EventCardBoundary from './ErrorBoundaries/EventCardBoundary';
import { useAuth } from './AuthProvider';
import EventCard from './EventCard';
import PoolSelector from './PoolSelector';

const EventList = () => {
    const user = useAuth();
    const { events, loading } = useAllEvents();
    const [displayPool, setDisplayPool] = useState<Pool>(Pool.FridayEvening);
    const [showFilters, setShowFilters] = useState(false);
    const [showUnpublished, setShowUnpublished] = useState(false);
    const { userSettings } = useUserSettings(user?.uid);
    const searchParams = useSearchParams();

    useEffect(() => {
        setShowUnpublished(userSettings?.admin && user ? true : false);
    }, [user, userSettings]);

    const search = searchParams.get('pool') as keyof typeof Pool;
    useEffect(() => {
        if (search) {
            setDisplayPool(Pool[search]);
        }
    }, [search]);

    return (
        <>
            <PoolSelector handlePoolChange={(pool) => setDisplayPool(pool)} />
            {/* {displayPool === 'Fredag Kveld' ? <Typography variant="h2">Fredag 17-22</Typography> : null}
            {displayPool === 'Lørdag Morgen' ? <Typography variant="h2">Lørdag 10-16</Typography> : null}
            {displayPool === 'Lørdag Kveld' ? <Typography variant="h2">Lørdag 17-22</Typography> : null}
            {displayPool === 'Søndag Morgen' ? <Typography variant="h2">Søndag 10-16</Typography> : null} */}

            <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20">
                {loading ? <Typography variant="body1">Loading...</Typography> : null}
       {/*          <Box sx={{ display: 'flex', gap: '.5em', flexGrow: '1', justifyContent: 'center', width: '100%' }}>
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
                    .toSorted((a, b) => a.sortingIndex - b.sortingIndex)
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
