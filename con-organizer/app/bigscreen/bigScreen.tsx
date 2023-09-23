'use client';

import { ErrorBoundary } from 'react-error-boundary';
import { Box, Card, CardContent, CardHeader } from '@mui/material';
import EventCardBoundary from '@/components/ErrorBoundaries/EventCardBoundary';
import EventCard from '@/components/EventCard';
import { Pool } from '@/lib/enums';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';

const BigScreen = () => {
    const { events, loading } = useAllEvents();
    if (loading) {
        return <h1>Loading...</h1>;
    }
    // throw new Error('test');

    return (
        <Box className="flex flex-row flex-wrap justify-center gap-4">
            <Box className="flex flex-col gap-4 bg-slate-800 p-4" sx={{ maxWidth: '440px' }}>
                <h1>Fredag Kveld</h1>
                <p>18:00 - 23:00</p>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Innsjekk Fredag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 16:30 - 18:00 </p>
                    </CardContent>
                </Card>
                {events
                    ?.filter((conEvent) => conEvent.pool === Pool.FridayEvening)
                    .map((conEvent) => (
                        <ErrorBoundary FallbackComponent={EventCardBoundary} key={conEvent.id}>
                            <EventCard conEvent={conEvent} />
                        </ErrorBoundary>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4 bg-slate-900 p-4" sx={{ maxWidth: '440px' }}>
                <h1>Lørdag Morgen </h1>
                <p>10:00 - 15:00</p>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Innsjekk Lørdag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 09:00 - 10:00 </p>
                    </CardContent>
                </Card>
                {events
                    ?.filter((conEvent) => conEvent.pool === Pool.SaturdayMorning)
                    .map((conEvent) => (
                        <ErrorBoundary FallbackComponent={EventCardBoundary} key={conEvent.id}>
                            <EventCard conEvent={conEvent} />
                        </ErrorBoundary>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4 bg-slate-800 p-4" sx={{ maxWidth: '440px' }}>
                <h1>Lørdag Kveld </h1>
                <p>18:00 - 23:00</p>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Middag Lørdag" />

                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 16:00 - 18:00 </p>
                        <p>Påmeling til middag er frem til 4. oktober for eksempel?</p>
                    </CardContent>
                </Card>
                {events
                    ?.filter((conEvent) => conEvent.pool === Pool.SaturdayEvening)
                    .map((conEvent) => (
                        <ErrorBoundary FallbackComponent={EventCardBoundary} key={conEvent.id}>
                            <EventCard conEvent={conEvent} />
                        </ErrorBoundary>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4 bg-slate-900 p-4" sx={{ maxWidth: '440px' }}>
                <h1>Søndag Morgen </h1>
                <p>10:00 - 15:00</p>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Innsjekk Søndag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 09:00 - 10:00 </p>
                    </CardContent>
                </Card>
                {events
                    ?.filter((conEvent) => conEvent.pool === Pool.SundayMorning)
                    .map((conEvent) => (
                        <ErrorBoundary FallbackComponent={EventCardBoundary} key={conEvent.id}>
                            <EventCard conEvent={conEvent} />
                        </ErrorBoundary>
                    ))}

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Prisutdeling" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 16:00 - 17:00 </p>
                    </CardContent>
                </Card>
            </Box>
        </Box>
    );
};
export default BigScreen;
