'use client';
import { Box, Typography } from '@mui/material';
import type { ConEvents, EventDay } from '../page';
import EventCardBig from './components/EventCardBig';
import EventCardSmall from './components/EventCardSmall';
import NextLink from 'next/link';
import { use, useEffect, useRef, useState } from 'react';
import EventListDay from './components/ui/EventListDay';
import { useIntersectionObserver } from './lib/hooks/useIntersectionObserver';
import { useObserveIntersectionObserver } from './lib/hooks/useObserveIntersectionObserver';
type Props = {
    events: ConEvents;
};

const EventList = ({ events }: Props) => {
    return (
        <Box>
            {events.map((event) => {
                const ref = useRef<HTMLDivElement>(null);
                useObserveIntersectionObserver(ref);

                return (
                    <Box key={event.day} ref={ref}>
                        <EventListDay eventDay={event.day} />
                        <Box
                            sx={{
                                display: 'grid',
                                gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                                gap: '1rem',
                            }}
                        >
                            {event.events.map((event) => (
                                <NextLink key={event.id} href={`/event/${event.id}`}>
                                    {event.isSmallCard ?
                                        <EventCardSmall
                                            title={event.title}
                                            gameMaster={event.gameMaster}
                                            system={event.system}
                                        />
                                        : <EventCardBig
                                            title={event.title}
                                            gameMaster={event.gameMaster}
                                            shortDescription={event.shortDescription}
                                            system={event.system}
                                            backgroundImage={event.smallImageURL}
                                        />
                                    }
                                </NextLink>
                            ))}
                        </Box>
                    </Box>
                );
            })}
        </Box>
    );
};

export default EventList;
