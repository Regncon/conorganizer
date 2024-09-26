'use client';
import type { IconName, PoolEvent } from '$lib/types';
import { Box } from '@mui/material';
import type { Route } from 'next';
import NextLink from 'next/link';
import EventCardBig from './EventCardBig';
import EventCardSmall from './EventCardSmall';
import { filteredEvents } from '../lib/helpers/filtersHelpers';

type Props = {
    events: PoolEvent[];
    searchParams: {
        [key in IconName]: string;
    };
};

const Events = ({ events, searchParams }: Props) => {
    // const filteredEvents = useFilteredEvents(events, searchParams);

    return filteredEvents(events, searchParams).map((event) => (
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
    ));
};

export default Events;
