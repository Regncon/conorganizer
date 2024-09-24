import { Box } from '@mui/material';
import EventCardBig from './components/EventCardBig';
import EventCardSmall from './components/EventCardSmall';
import NextLink from 'next/link';
import EventListDay from './components/ui/EventListDay';
import type { PoolEvents } from './lib/serverAction';
import { Route } from 'next';
import type { PoolEvent } from '$lib/types';
import EventListWrapper from './components/EventListWrapper';

type Props = {
    events: PoolEvents;
};

const removeUnpublishedEventsFilter = (event: PoolEvent) => event.published;
const EventList = ({ events }: Props) => {
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
                            {events.filter(removeUnpublishedEventsFilter).map((event) => (
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
                    </EventListWrapper>
                );
            })}
        </Box>
    );
};

export default EventList;
