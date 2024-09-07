'use client';
import { Box, Typography } from '@mui/material';
import type { ConEvents } from '../../page';
import EventCardBig from '../components/EventCardBig';
import EventCardSmall from '../components/EventCardSmall';
import NextLink from 'next/link';
import { use, useEffect, useRef } from 'react';
import EventListDay from './EventListDay';
import { IntersectionObserverContext } from '../lib/IntersectionObserverContext';
type Props = {
    events: ConEvents;
};

const EventList = ({ events }: Props) => {
    return events.map((event) => {
        const ref = useRef<HTMLDivElement>(null);
        const intersectionObserver = use(IntersectionObserverContext);

        useEffect(() => {
            if (ref.current) {
                intersectionObserver?.observe(ref.current);
            }
            return () => {
                if (ref.current) {
                    intersectionObserver?.unobserve(ref.current);
                }
            };
        }, [ref, ref.current]);
        return (
            <Box key={event.day} ref={ref}>
                <EventListDay eventDay={event.day} />
                <Box sx={{ display: 'grid' }}>
                    {event.events.map((event) => (
                        <NextLink key={event.id} href={`/event/${event.id}`} style={{ all: 'unset' }}>
                            {event.isSmallCard ?
                                <EventCardSmall
                                    title={event.title}
                                    gameMaster={event.gameMaster}
                                    system={event.system}
                                />
                            :   <EventCardBig
                                    title={event.title}
                                    gameMaster={event.gameMaster}
                                    shortDescription={event.shortDescription}
                                    system={event.system}
                                />
                            }
                        </NextLink>
                    ))}
                </Box>
            </Box>
        );
    });
};

export default EventList;
