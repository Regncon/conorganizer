import EventCardBig from './EventCardBig';
import EventCardSmall from './EventCardSmall';
import { getAllEvents } from './serverAction';
import RealtimeEvents from './RealtimeEvents';
import Grid from '@mui/material/Unstable_Grid2';
import { redirect } from 'next/navigation';
import { redirectToAdminDashboardWhenAdministrator } from '$lib/lib';

export default async function Home() {
    await redirectToAdminDashboardWhenAdministrator();

    redirect('/dashboard');
    const events = await getAllEvents();

    return (
        <>
            <img src="/placeholderlogo.png" alt="logo" />
            <Grid container spacing={2}>
                {events.map((event, i) => {
                    return (
                        <Grid xs={i === 0 ? 12 : 6}>
                            {i === 0 ?
                                <Grid xs={12}>
                                    <EventCardBig
                                        key={i}
                                        title={event.title}
                                        gameMaster={event.gameMaster}
                                        shortDescription={event.shortDescription}
                                        system={event.system}
                                    />
                                </Grid>
                                : <EventCardSmall
                                    key={i}
                                    title={event.title}
                                    gameMaster={event.gameMaster}
                                    system={event.system}
                                />
                            }
                        </Grid>
                    );
                })}
            </Grid>
            {/* <EventCardBig
    					title="Hello world"
    					gameMaster="Gerhard Fajita"
    					shortDescription="Mord overalt! Kos! Gøy!"
    					system="Call of Chthuhlth"
    				/>
    				<Box sx={{ display: 'flex' }}>
    					<EventCardSmall title="Hi" gameMaster="Gardh Fajita2" system="Dungeons 2" />
    					<EventCardSmall title="Any% speedrun" gameMaster="Gorde Fajita3" system="Terraria" />
    				</Box> */}

            <RealtimeEvents />
        </>
    );
}
