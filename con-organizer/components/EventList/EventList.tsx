'use client';

import { useEffect, useState } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { useSearchParams } from 'next/navigation';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { useUserSettings } from '@/lib/hooks/UseUserSettings';
import { CustomEventTypeNames, Pool } from '@/models/enums';
import { ConEvent, CustomEventTypeFilteredEvent } from '@/models/types';
import { Box, Typography } from '../../lib/mui';
import { useAuth } from '../AuthProvider/AuthProvider';
import EventCardBoundary from '../ErrorBoundaries/EventCardBoundary';
import EventCard from '../Event/EventCard';
import PoolSelector from '../Navigation/PoolSelector';
import Filters from './Filters/Filter';

const EventList = () => {
    const user = useAuth();
    const { loading } = useAllEvents();
    const [showUnpublished, setShowUnpublished] = useState(false);
    const [filteredEvents, setFilteredEvents] = useState<ConEvent[] | undefined>();
    const [displayPool, setDisplayPool] = useState<Pool>(Pool.FridayEvening);
    const searchParams = useSearchParams();
    const { userSettings } = useUserSettings(user?.uid);
    useEffect(() => {
        setShowUnpublished(userSettings?.admin && user ? true : false);
    }, [user, userSettings]);

    const search = searchParams.get('pool') as keyof typeof Pool;
    const handleUpdateFilteredChanges = (e: CustomEvent<CustomEventTypeFilteredEvent>) => {
        setFilteredEvents(e.detail.filteredEvents);
    };
    useEffect(() => {
        window.addEventListener(CustomEventTypeNames.FilterChanges, handleUpdateFilteredChanges);
        return () => {
            window.removeEventListener(CustomEventTypeNames.FilterChanges, handleUpdateFilteredChanges);
        };
    }, []);

    useEffect(() => {
        if (search) {
            setDisplayPool(Pool[search]);
        }
    }, [search]);

    return (
        <>
            <PoolSelector handlePoolChange={(pool) => setDisplayPool(pool)} />
            <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20">
                {loading ? <Typography variant="body1">Loading...</Typography> : null}
                <Filters displayPool={displayPool} />
                {filteredEvents
                    ?.filter((conEvent) => showUnpublished || conEvent.published)
                    .sort((a, b) => a.sortingIndex - b.sortingIndex)
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
