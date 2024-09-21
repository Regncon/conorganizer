'use client';
import { Box } from '@mui/material';
import EventCardBig from './components/EventCardBig';
import EventCardSmall from './components/EventCardSmall';
import NextLink from 'next/link';
import { useRef } from 'react';
import EventListDay from './components/ui/EventListDay';
import { useObserveIntersectionObserver } from './lib/hooks/useObserveIntersectionObserver';
import type { PoolEvents } from './lib/serverAction';
import { Route } from 'next';

type Props = {
    events: PoolEvents;
};

const EventList = ({ events }: Props) => {
    return (
        <Box>
            {[...events.entries()].map(([day, events]) => {
                const ref = useRef<HTMLDivElement>(null);
                useObserveIntersectionObserver(ref);
                return (
                    <Box key={day} ref={ref}>
                        <EventListDay poolDay={day} />
                        <Box
                            sx={{
                                display: 'grid',
                                gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                                gap: '1rem',
                            }}
                        >
                            {events.map((event) => (
                                <Box
                                    component={NextLink}
                                    key={event.id}
                                    sx={{ textDecoration: 'none' }}
                                    prefetch
                                    href={`/event/${event.id}` as Route}
                                >
                                    {event.isSmallCard ?
                                        <EventCardSmall
                                            title={event.title}
                                            gameMaster={event.gameMaster}
                                            system={event.system}
                                            backgroundImage={event.smallImageURL}
                                        />
                                    :   <EventCardBig
                                            title={event.title}
                                            gameMaster={event.gameMaster}
                                            shortDescription={event.shortDescription}
                                            system={event.system}
                                            backgroundImage={event.smallImageURL}
                                        />
                                    }
                                </Box>
                            ))}
                        </Box>
                    </Box>
                );
            })}
        </Box>
    );
};

export default EventList;
