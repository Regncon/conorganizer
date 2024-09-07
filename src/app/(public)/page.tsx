import Image from 'next/image';
import { Box } from '@mui/material';

import RealtimeEvents from './components/RealtimeEvents';
import { getAllEvents } from './components/lib/serverAction';
import DaysHeader from './components/ui/DaysHeader';
import type { ConEvent } from '$lib/types';
import EventList from './components/EventList';
import Logo from './components/ui/Logo';

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
                <Logo />
                <DaysHeader eventDays={eventDays} />
                <EventList events={events} />
            </Box>
            <RealtimeEvents where="EVENTS" />
        </>
    );
}
