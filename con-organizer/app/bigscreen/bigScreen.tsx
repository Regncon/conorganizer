'use client';

import { Box, Card, CardContent, CardHeader } from '@mui/material';
import EventHeader from '@/components/eventHeader';
import { Pool } from '@/lib/enums';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';

const BigScreen = () => {
    const { events, loading } = useAllEvents();
    console.log(events);

    if (loading) {
        return <h1>Loading...</h1>;
    }

    return (
        <Box className="flex flex-row flex-wrap justify-center gap-4">
            <Box className="flex flex-col gap-4">
                <h1>Fredag Kveld </h1>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Registrering Fredag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 16:00 - 17:00 </p>
                    </CardContent>
                </Card>
                {events
                    ?.filter((conEvent) => conEvent.pool === Pool.FridayEvening)
                    .map((conEvent) => (
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4">
                <h1>Lørdag Morgen </h1>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Registrering Lørdag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 10:00 - 11:00 </p>
                    </CardContent>
                </Card>
                {events
                    ?.filter((conEvent) => conEvent.pool === Pool.SaturdayMorning)
                    .map((conEvent) => (
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4">
                <h1>Lørdag Kveld </h1>

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
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    ))}
            </Box>

            <Box className="flex flex-col gap-4">
                <h1>Søndag Morgen </h1>

                <Card sx={{ width: '100%' }}>
                    <CardHeader sx={{ paddingBottom: '0.5rem' }} title="Registrering Søndag" />
                    <CardContent sx={{ paddingTop: '0' }}>
                        <p>Kl 10:00 - 11:00 </p>
                    </CardContent>
                </Card>
                {events
                    ?.filter((conEvent) => conEvent.pool === Pool.SundayMorning)
                    .map((conEvent) => (
                        <Card
                            key={conEvent.id}
                            onClick={() => {
                                window.location.assign(`/event/${conEvent.id}`);
                            }}
                            sx={{ cursor: 'pointer', opacity: conEvent?.published === false ? '50%' : '' }}
                        >
                            <EventHeader conEvent={conEvent} />
                        </Card>
                    ))}
            </Box>
        </Box>
    );
};
export default BigScreen;
