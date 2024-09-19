import { Grid2, Paper, Typography } from '@mui/material';
import EventForm from './components/EventForm';
import type { Metadata } from 'next';
import { getMyEventById } from '$app/(public)/components/lib/serverAction';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';

type Props = {
    params: {
        id: string;
    };
};
export async function generateMetadata({ params: { id } }: Props): Promise<Metadata> {
    const { user } = await getAuthorizedAuth();

    if (user) {
        const event = await getMyEventById(id, user.uid);
        if (event && event.title && event.title.length > 0) {
            return {
                title: `Oppdater arrangementet ${event.title}`,
            };
        }
    }
    return {};
}
const Create = async ({ params: { id } }: Props) => {
    return (
        <>
            <Paper sx={{ padding: '1rem', marginBottom: '2rem' }}>
                <Grid2 container gap="3rem">
                    <Typography variant="h1">Meld på arrangement til Regncon XXXII 2024</Typography>
                    <Typography>
                        Takk for at du vil arrangere eit spel på Regncon, anten det er brettspel, kortspel, rollespel
                        eller anna, så sett vi enormt pris på ditt bidrag. Fyll inn skjemaet så godt du kan, og ikkje
                        vere redd for å ta kontakt med Regnconstyret på{' '}
                        <a href="mailto:regncon@gmail.com">regncon@gmail.com</a> om du skulle lure på noko!
                    </Typography>
                </Grid2>
            </Paper>
            <EventForm id={id} />
        </>
    );
};
export default Create;
