import { Box } from '@mui/material';
import EventCardBig from './components/EventCardBig';
import EventCardSmall from './components/EventCardSmall';
import NextLink from 'next/link';
import EventListDay from './components/ui/EventListDay';
import type { PoolEvents } from './lib/serverAction';
import { Route } from 'next';
import type { IconName, PoolEvent } from '$lib/types';
import EventListWrapper from './components/EventListWrapper';
import Events from './components/Events';

type Props = {
    events: PoolEvents;
    searchParams: {
        [key in IconName]: string;
    };
};

const EventList = ({ events, searchParams }: Props) => {
    return (
        <Box>
            {[...events.entries()].map(([day, events]) => {
                return (
                    <EventListWrapper key={day} day={day}>
                        <EventListDay poolDay={day} />
                        <Box
                            sx={{
                                display: 'grid',
                                gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                                gap: '1rem',
                            }}
                        >
                            {/* <Events events={events} searchParams={searchParams} /> */}
                        </Box>
                    </EventListWrapper>
                );
            })}
        </Box>
    );
};

export default EventList;
