import EventCardBig from '$app/(public)/EventCardBig';
import { getAllEvents } from '$app/(public)/serverAction';
import { Grid2, Paper } from '@mui/material';
import Link from 'next/link';

const Events = async () => {
    const events = await getAllEvents();

    return (
        <Paper elevation={0}>
            <Grid2 gap={'2rem'} container sx={{ padding: '2rem' }}>
                {events.map((event) => {
                    return (
                        <Link
                            href={`/admin/dashboard/events/event-dashboard/${event.id}/edit`}
                            style={{ textDecoration: 'none' }}
                            key={event.id}
                        >
                            <EventCardBig
                                title={event.title}
                                gameMaster={event.gameMaster}
                                system={event.system}
                                shortDescription={event.shortDescription}
                            />
                        </Link>
                    );
                })}
            </Grid2>
        </Paper>
    );
};

export default Events;
