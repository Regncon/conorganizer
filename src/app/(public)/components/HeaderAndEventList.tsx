'use client';

import { Box, Typography } from '@mui/material';
import EventCardBig from './components/EventCardBig';
import EventCardSmall from './components/EventCardSmall';
import type { ConEvents, EventDay, EventDays } from '../page';

import { useCallback, useEffect, useRef, useState } from 'react';
import DaysHeader from './ui/DaysHeader';
import EventList from './ui/EventList';
import debounce from '$lib/debounce';

type Props = {
    events: ConEvents;
    eventDays: EventDays;
};
const HeaderAndEventList = ({ events, eventDays }: Props) => {
    const ref = useRef<HTMLDivElement>(null);
    const [locationHash, setLocationHash] = useState<EventDay>('');
    const [intersectionObserver, setIntersectionObserver] = useState<IntersectionObserver | null>(null);

    useEffect(() => {
        if (window) {
            setLocationHash((prev) => decodeURI(window.location.hash) as EventDay);
        }

        const handleIntersectionObserver: (entries: IntersectionObserverEntry[]) => void = (entries) => {
            entries.forEach(async (entry) => {
                if (entry.isIntersecting) {
                    debounce(() => {
                        if (decodeURI(window.location.hash) === `#${entry.target.id}`) {
                            const id = entry.target.id as EventDay;
                            setLocationHash(id);
                        }
                    }, 200)();
                }
            });
        };

        setIntersectionObserver(
            new IntersectionObserver(handleIntersectionObserver, {
                root: null,
                threshold: 1,
            })
        );
    }, []);

    return (
        <>
            <DaysHeader eventDays={eventDays} locationHash={locationHash} />
            <Box ref={ref}>
                <EventList events={events} intersectionObserver={intersectionObserver} />
            </Box>
        </>
    );
};

export default HeaderAndEventList;
