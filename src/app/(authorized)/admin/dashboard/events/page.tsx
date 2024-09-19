import EventCardBig from '$app/(public)/components/components/EventCardBig';
import RealtimeEvents from '$app/(public)/components/RealtimeEvents';
import { getAllEvents } from '$app/(public)/components/lib/serverAction';
import { Box, Paper } from '@mui/material';
import Link from 'next/link';
import type { Metadata } from 'next';
export const metadata: Metadata = {
    title: 'Liste over arrangementer som kan administreres',
};

const Events = async () => {
    const events = await getAllEvents();
    return (
        <>
            <Paper elevation={0}>
                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                        gap: '1rem',
                    }}
                >
                    {events
                        .sort((a, b) => a.title.localeCompare(b.title))
                        .map((conEvent) => {
                            return (
                                <Link
                                    href={`/admin/dashboard/events/event-dashboard/${conEvent.id}/edit`}
                                    prefetch
                                    style={{ textDecoration: 'none' }}
                                    key={conEvent.id}
                                >
                                    <EventCardBig
                                        title={conEvent.title}
                                        gameMaster={conEvent.gameMaster}
                                        system={conEvent.system}
                                        shortDescription={conEvent.shortDescription}
                                        backgroundImage={conEvent.smallImageURL}
                                    />
                                </Link>
                            );
                        })}
                </Box>
            </Paper>
            <RealtimeEvents where="DASHBOARD_EVENTS" />
        </>
    );
};

export default Events;
