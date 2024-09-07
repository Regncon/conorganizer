'use client';

import { Box, Typography } from '@mui/material';
import EventCardBig from './components/EventCardBig';
import EventCardSmall from './components/EventCardSmall';
import type { ConEvents, EventDay, EventDays } from '../page';

import { useCallback, useEffect, useRef, useState, useTransition } from 'react';
import DaysHeader from './ui/DaysHeader';
import EventList from './ui/EventList';
import debounce from '$lib/debounce';
import { IntersectionObserverContext } from './lib/IntersectionObserverContext';

type Props = {
    events: ConEvents;
    eventDays: EventDays;
};
const HeaderAndEventList = ({ events, eventDays }: Props) => {
    const ref = useRef<HTMLDivElement>(null);
    const [locationHash, setLocationHash] = useState<EventDay>('');
    const [intersectionObserver, setIntersectionObserver] = useState<IntersectionObserver | null>(null);

    useEffect(() => {
        if (typeof window !== 'undefined' && typeof document !== 'undefined') {
            setLocationHash(decodeURI(window.location.hash).substring(1) as EventDay);

            const handleIntersectionObserver: (entries: IntersectionObserverEntry[]) => void = (entries) => {
                entries.forEach(async (entry) => {
                    if (entry.isIntersecting) {
                        const id = entry.target.querySelector('h1')?.id as EventDay;
                        debounce(() => {
                            setLocationHash(id);
                        }, 450)();
                    }
                });
            };

            setIntersectionObserver(
                new IntersectionObserver(handleIntersectionObserver, {
                    root: null,
                    threshold: 0,
                })
            );
        }
    }, []);

    return (
        <IntersectionObserverContext.Provider value={intersectionObserver}>
            {intersectionObserver ?
                <DaysHeader eventDays={eventDays} locationHash={locationHash} />
            :   null}
            <Box ref={ref}>
                <EventList events={events} />
            </Box>
        </IntersectionObserverContext.Provider>
    );
};

export default HeaderAndEventList;
