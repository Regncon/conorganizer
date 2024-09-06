import EventCardBig from './components/EventCardBig';
import Image from 'next/image';
import { Box, Link, Typography } from '@mui/material';
import NextLink from 'next/link';

import EventCardSmall from './components/EventCardSmall';
import RealtimeEvents from './components/RealtimeEvents';
import { getAllEvents } from './components/serverAction';
import DaysHeader from './components/ui/DaysHeader';

const eventDays = {
    fridayEvening: 'Fredag',
    saturdayMorning: 'Lørdag Morgen',
    saturdayEvening: 'Lørdag Kveld',
    sunday: 'Søndag',
};
export type EventDays = typeof eventDays;
export default async function Home() {
    const allEvents = await getAllEvents();

    const events = [
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

                <DaysHeader eventDays={eventDays} />

                <Box>
                    {events.map((event) => {
                        return (
                            <Box key={event.day}>
                                <Typography
                                    id={event.day}
                                    sx={{ scrollMarginTop: 'var(--scroll-margin-top)' }}
                                    variant="h1"
                                >
                                    {event.day}
                                </Typography>
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
                    })}
                </Box>
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
