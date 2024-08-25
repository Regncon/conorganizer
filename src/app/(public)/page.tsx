import EventCardBig from './EventCardBig';
import EventCardSmall from './EventCardSmall';
import { getAllEvents } from './serverAction';
import RealtimeEvents from './RealtimeEvents';
import Grid from '@mui/material/Unstable_Grid2';
import Image from 'next/image';
import { Box } from '@mui/material';
import Link from 'next/link';

export default async function Home() {
    const events = await getAllEvents();

    return (
        <Box>
            <Link href="/">
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
            </Link>
            <Grid container spacing={0}>
                <Grid container spacing={0}>
                    {events
                        .filter((ce) => ce.published)
                        .map((event, i) => {
                            return event.isSmallCard ?
                                <Grid xs={6}>
                                    <Grid
                                        xs
                                        display="flex"
                                        justifyContent="center"
                                        alignItems="center"
                                        paddingBottom={'1rem'}
                                    >
                                        <EventCardSmall
                                            key={i}
                                            title={event.title}
                                            gameMaster={event.gameMaster}
                                            system={event.system}
                                        />
                                    </Grid>
                                </Grid>
                                : <Grid xs={12}>
                                    <Grid
                                        xs
                                        display="flex"
                                        justifyContent="center"
                                        alignItems="center"
                                        paddingBottom={'1rem'}
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
                                    </Grid>
                                </Grid>;
                        })}
                </Grid>
                <RealtimeEvents />
            </Grid>
        </Box>
    );
}
