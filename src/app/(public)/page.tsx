import EventCardBig from './components/components/EventCardBig';
import Image from 'next/image';
import { Box, Link, Typography } from '@mui/material';
import NextLink from 'next/link';

import EventCardSmall from './components/components/EventCardSmall';
import RealtimeEvents from './components/RealtimeEvents';
import { getAllEvents } from './components/lib/serverAction';
import DaysHeader from './components/ui/DaysHeader';
import type { ConEvent } from '$lib/types';
import EventsList from './components/HeaderAndEventList';

export type EventDays = typeof eventDays;
export type EventDay = EventDays[keyof EventDays] | '';
export type ConEvents = {
    day: EventDay;
    events: ConEvent[];
}[];
const eventDays = {
    fridayEvening: 'Fredag',
    saturdayMorning: 'Lørdag Morgen',
    saturdayEvening: 'Lørdag Kveld',
    sunday: 'Søndag',
} as const;
export default async function Home() {
    const allEvents = await getAllEvents();

    const events: ConEvents = [
        { day: eventDays.fridayEvening, events: [...allEvents] },
        { day: eventDays.saturdayMorning, events: [...allEvents] },
        { day: eventDays.saturdayEvening, events: [...allEvents] },
        { day: eventDays.sunday, events: [...allEvents] },
    ];

    return (
        <>
            <Box>
                <Box
                    sx={{
                        maxWidth: '430px',
                        maxHeight: '430px',
                        margin: 'auto',
                        width: '100vw',
                        aspectRatio: '1/1',
                        marginBlockStart: '0.5rem',
                        marginBlockEnd: '1rem',
                        position: 'relative',
                    }}
                >
                    <Image src="/RegnCon2024LogoWhite.webp" fill alt="logo" />
                </Box>

                <EventsList events={events} eventDays={eventDays} />

                {/* <Grid2 container spacing={0}>
                    <Grid2 container spacing={0}>
                        {events
                            .filter((ce) => ce.published)
                            .map((event, i) => {
                                return event.isSmallCard ?
                                        <Grid2 size={6}>
                                            <Grid2
                                                display="flex"
                                                justifyContent="center"
                                                alignItems="center"
                                                paddingBottom={'1rem'}
                                                size="grow"
                                            >
                                                <EventCardSmall
                                                    key={i}
                                                    title={event.title}
                                                    gameMaster={event.gameMaster}
                                                    system={event.system}
                                                />
                                            </Grid2>
                                        </Grid2>
                                    :   <Grid2 size={12}>
                                            <Grid2
                                                display="flex"
                                                justifyContent="center"
                                                alignItems="center"
                                                paddingBottom={'1rem'}
                                                size="grow"
                                            >
                                                <Link href={`/event/${event.id}`} style={{ all: 'unset' }}>
                                                    <EventCardBig
                                                        key={i}
                                                        title={event.title}
                                                        gameMaster={event.gameMaster}
                                                        shortDescription={event.shortDescription}
                                                        system={event.system}
                                                    />
                                                </Link>
                                            </Grid2>
                                        </Grid2>;
                            })}
                    </Grid2>
                </Grid2> */}
            </Box>
            <RealtimeEvents where="EVENTS" />
        </>
    );
}
