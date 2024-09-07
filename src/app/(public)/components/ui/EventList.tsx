'use client';
import { Box, Typography } from '@mui/material';
import type { ConEvents } from '../../page';
import EventCardBig from '../components/EventCardBig';
import EventCardSmall from '../components/EventCardSmall';
import NextLink from 'next/link';
import { useEffect, useRef } from 'react';
import EventListDay from './EventListDay';
type Props = {
    events: ConEvents;
    intersectionObserver: IntersectionObserver | null;
};

const EventList = ({ events, intersectionObserver }: Props) => {
    const ref = useRef<HTMLDivElement>(null);
    // const intersectionObserver = new IntersectionObserver(
    //     (entries) => {
    //         entries.forEach((entry) => {
    //             console.log(entry);
    //         });
    //     },
    //     { root: ref.current, threshold: 1 }
    // );

    // useEffect(() => {
    //     if (ref.current) {
    //         console.log('test');

    //         intersectionObserver.observe(ref.current);
    //     }
    //     return () => {
    //         if (ref.current) {
    //             intersectionObserver.unobserve(ref.current);
    //         }
    //     };
    // }, [ref, ref.current]);
    return events.map((event) => {
        return (
            <Box key={event.day}>
                <EventListDay eventDay={event.day} intersectionObserver={intersectionObserver} />
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
