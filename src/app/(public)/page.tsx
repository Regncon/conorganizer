import EventCardBig from './EventCardBig';
import EventCardSmall from './EventCardSmall';
import { getAllEvents } from './serverAction';
import RealtimeEvents from './RealtimeEvents';
import Grid from '@mui/material/Unstable_Grid2';
import Image from 'next/image';
import { Box } from '@mui/material';

export default async function Home() {
    const events = await getAllEvents();

    return (
        <Box sx={{ display: 'grid', placeItems: 'center', placeContent:'center' }}>
            <Box
                sx={{
                    width: { md: '300px', sm: '250px', xs: '200px' },
                    height: { md: '300px', sm: '250px', xs: '200px' },
                    marginBlockStart: '0.5rem',
                    marginBlockEnd: '1rem',
                    position: 'relative',
                }}
            >
                <Image src="/RegnCon2024LogoWhite.webp" fill alt="logo" />
            </Box>
            <Grid container spacing={2}>
                <Grid container spacing={2}>
                    {events
                        .filter((ce) => ce.published)
                        .map((event, i) => {
                            return event.isSmallCard ?
                                    <Grid xs={6}>
                                        <EventCardSmall
                                            key={i}
                                            title={event.title}
                                            gameMaster={event.gameMaster}
                                            system={event.system}
                                        />
                                    </Grid>
                                :   <Grid xs={12}>
                                        <EventCardBig
                                            key={i}
                                            title={event.title}
                                            gameMaster={event.gameMaster}
                                            shortDescription={event.shortDescription}
                                            system={event.system}
                                        />
                                    </Grid>;
                        })}
                </Grid>
                <RealtimeEvents />
            </Grid>
        </Box>
    );
}
